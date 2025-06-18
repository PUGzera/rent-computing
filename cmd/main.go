package main

import (
	"context"
	"fmt"
	"log"
	"os"
	user_controller "rent-computing/internal/user/controller"
	user_repo "rent-computing/internal/user/repo"
	user_rest_handler "rent-computing/internal/user/rest_handler"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	userRepo, err := user_repo.NewMongoDB(user_repo.OptionsMongoDB{
		Collection: db.Collection("users"),
	})
	if err != nil {
		log.Fatalf("failed to create user repository: %s", err.Error())
	}

	userController, err := user_controller.New(user_controller.Options{
		UserRepo: userRepo,
	})
	if err != nil {
		log.Fatalf("failed to create user controller: %s", err.Error())
	}

	_, err = user_rest_handler.New(user_rest_handler.Options{
		Router:         router.Group("/users"),
		JWTSecret:      "secret",
		UserController: userController,
	})
	if err != nil {
		log.Fatalf("failed to create user rest handler: %s", err.Error())
	}

	router.Run()
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
