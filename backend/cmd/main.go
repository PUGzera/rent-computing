package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	machine_controller "rent-computing/internal/machine/controller"
	machine_manager "rent-computing/internal/machine/manager"
	machine_repo "rent-computing/internal/machine/repo"
	machine_rest_handler "rent-computing/internal/machine/rest_handler"
	user_controller "rent-computing/internal/user/controller"
	user_repo "rent-computing/internal/user/repo"
	user_rest_handler "rent-computing/internal/user/rest_handler"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	conn, err := connectMongoDB(time.Minute * 2)
	if err != nil {
		log.Fatalf("failed to create a psql connection: %s", err.Error())
	}

	dbName, ok := os.LookupEnv("MONGODB_DB")
	if !ok {
		dbName = "rent-computing"
	}
	db := conn.Database(dbName)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userRepo, err := user_repo.NewMongoDB(user_repo.OptionsMongoDB{
		Collection: db.Collection("users"),
	})
	if err != nil {
		log.Fatalf("failed to create user repository: %s", err.Error())
	}
	machineRepo, err := machine_repo.NewMongoDB(machine_repo.OptionsMongoDB{
		Collection: db.Collection("machines"),
	})
	if err != nil {
		log.Fatalf("failed to create machine repository: %s", err.Error())
	}

	userController, err := user_controller.New(user_controller.Options{
		UserRepo: userRepo,
	})
	if err != nil {
		log.Fatalf("failed to create user controller: %s", err.Error())
	}

	kubeconfig, err := kubeconfig()
	if err != nil {
		log.Fatalf("failed to get kubeconfig: %s", err.Error())
	}

	namespace, ok := os.LookupEnv("NAMESPACE")
	if !ok {
		namespace = "default"
	}

	machineRoute := "/machines"

	machineManager, err := machine_manager.NewK8s(machine_manager.OptionsK8s{
		Namespace: namespace,
		Config:    kubeconfig,
		Route:     machineRoute,
	})
	if err != nil {
		log.Fatalf("failed to create machine manager: %s", err.Error())
	}

	machineController, err := machine_controller.New(machine_controller.Options{
		MachineRepo:    machineRepo,
		MachineManager: machineManager,
	})
	if err != nil {
		log.Fatalf("failed to create machine controller: %s", err.Error())
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		jwtSecret = "secret"
	}

	jwtHandler, err := jwt.New(&jwt.GinJWTMiddleware{
		Key: []byte(jwtSecret),
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var userReq map[string]string
			if err := c.ShouldBindJSON(&userReq); err != nil {
				return nil, err
			}

			username, ok := userReq["username"]
			if !ok {
				return nil, errors.New("no username found in request")
			}

			password, ok := userReq["password"]
			if !ok {
				return nil, errors.New("no password found in request")
			}

			user, err := userController.Login(c, username, password)
			if err != nil {
				return nil, err
			}

			userReq["id"] = user.Id
			userReq["email"] = user.Email

			return userReq, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			ret := jwt.MapClaims{}
			if user, ok := data.(map[string]string); ok {
				ret["username"] = user["username"]
				ret["email"] = user["email"]
				ret["id"] = user["id"]
				return ret
			}
			return ret
		},
		Timeout:        time.Hour,
		MaxRefresh:     time.Hour,
		TokenLookup:    "cookie:jwt_token",
		SendCookie:     true,
		CookieName:     "jwt_token",
		CookieHTTPOnly: true,
		SecureCookie:   false,
		CookieSameSite: http.SameSiteStrictMode,
		CookieDomain:   "localhost",
	})
	if err != nil {
		log.Fatalf("failed to create jwt handler: %s", err.Error())
	}

	_, err = user_rest_handler.New(user_rest_handler.Options{
		Router:         router.Group("/users"),
		JWTHandler:     jwtHandler,
		UserController: userController,
	})
	if err != nil {
		log.Fatalf("failed to create user rest handler: %s", err.Error())
	}

	_, err = machine_rest_handler.New(machine_rest_handler.Options{
		Router:            router.Group(machineRoute),
		JWTHandler:        jwtHandler,
		MachineController: machineController,
		UserController:    userController,
	})
	if err != nil {
		log.Fatalf("failed to create user rest handler: %s", err.Error())
	}

	router.Run()
}

func kubeconfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	kubeconfigPath, ok := os.LookupEnv("KUBECONFIG")
	if ok {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err == nil {
			return config, nil
		}
	}

	home := homedir.HomeDir()
	kubeconfigPath = filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func connectMongoDB(timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return mongo.Connect(ctx, options.Client().ApplyURI(mongoConnString()))
}

func connectPSQL(timeout time.Duration) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fmt.Println(psqlConnString())
	return pgx.Connect(ctx, psqlConnString())
}

func mongoConnString() string {
	return connString("mongodb", "27017")
}

func psqlConnString() string {
	return connString("postgres", "5432")
}

func connString(dbType, defaultPort string) string {
	connString := fmt.Sprintf("%s://", dbType)

	username, ok := os.LookupEnv(fmt.Sprintf("%s_USER", strings.ToUpper(dbType)))
	if ok {
		connString = connString + username
		password, ok := os.LookupEnv(fmt.Sprintf("%s_PASS", strings.ToUpper(dbType)))
		if ok {
			connString = connString + ":" + password
		}
		connString = connString + "@"
	}

	address, ok := os.LookupEnv(fmt.Sprintf("%s_ADDR", strings.ToUpper(dbType)))
	if !ok || strings.EqualFold(address, "") {
		address = "localhost"
	}
	connString = connString + address

	port, ok := os.LookupEnv(fmt.Sprintf("%s_PORT", strings.ToUpper(dbType)))
	if !ok || strings.EqualFold(port, "") {
		port = defaultPort
	}
	connString = connString + ":" + port

	db, ok := os.LookupEnv(fmt.Sprintf("%s_DB", strings.ToUpper(dbType)))
	if !ok || strings.EqualFold(db, "") {
		db = "rent-computing"
	}
	connString = connString + "/" + db

	return connString
}
