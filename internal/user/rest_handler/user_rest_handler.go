package user_rest_handler

import (
	"errors"
	user_controller "rent-computing/internal/user/controller"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type createUserRequestBody struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type authenticateUserRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authenticatedUser struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type userRestHandler struct {
	router         *gin.RouterGroup
	jwtHandler     *jwt.GinJWTMiddleware
	userController user_controller.UserController
}

type Options struct {
	Router         *gin.RouterGroup
	JWTSecret      string
	UserController user_controller.UserController
}

func New(options Options) (*userRestHandler, error) {
	userController := options.UserController

	jwtHandler, err := jwt.New(&jwt.GinJWTMiddleware{
		Key: []byte(options.JWTSecret),
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var userReq authenticateUserRequestBody
			if err := c.ShouldBindJSON(&userReq); err != nil {
				return nil, err
			}

			user, err := userController.Login(c, userReq.Username, userReq.Password)
			if err != nil {
				return nil, err
			}

			return &authenticatedUser{
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			}, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			ret := jwt.MapClaims{}
			if user, ok := data.(*authenticatedUser); ok {
				ret["username"] = user.Username
				ret["email"] = user.Email
				return ret
			}
			return ret
		},
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
	})
	if err != nil {
		return nil, err
	}

	userRestHandler := userRestHandler{
		router:         options.Router,
		jwtHandler:     jwtHandler,
		userController: userController,
	}
	userRestHandler.registerUserRoute()
	userRestHandler.authenticateUserRoute()
	userRestHandler.getAuthenticatedUserRoute()
	userRestHandler.logoutUserRoute()
	userRestHandler.refreshUserRoute()

	return &userRestHandler, nil
}

func (r *userRestHandler) registerUserRoute() {
	r.router.POST("/register", func(ctx *gin.Context) {
		var newUser createUserRequestBody
		err := ctx.ShouldBindJSON(&newUser)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user, err := r.userController.Register(ctx, newUser.Username, newUser.Email, newUser.Password)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"user": authenticatedUser{
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}})
	})
}

func (r *userRestHandler) authenticateUserRoute() {
	r.router.POST("/login", r.jwtHandler.LoginHandler)
}

func (r *userRestHandler) logoutUserRoute() {
	r.router.
		Use(r.AuthenticationMiddleware()).
		GET("/logout", r.jwtHandler.LogoutHandler)
}

func (r *userRestHandler) refreshUserRoute() {
	r.router.
		Use(r.AuthenticationMiddleware()).
		GET("/refresh", r.jwtHandler.RefreshHandler)
}

func (r *userRestHandler) getAuthenticatedUserRoute() {
	r.router.GET("/profile", func(ctx *gin.Context) {
		user, err := r.GetAuthenticatedUser(ctx)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"user": user})
	})
}

func (r *userRestHandler) GetAuthenticatedUser(ctx *gin.Context) (*authenticatedUser, error) {
	claims, err := r.jwtHandler.GetClaimsFromJWT(ctx)
	if err != nil {
		return nil, err
	}

	usernameClaim, ok := claims["username"]
	if !ok {
		return nil, errors.New("claim \"username\" not found")
	}

	username, ok := usernameClaim.(string)
	if !ok {
		return nil, errors.New("username was not of type string")
	}

	user, err := r.userController.ProfileFromUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &authenticatedUser{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (r *userRestHandler) AuthenticationMiddleware() gin.HandlerFunc {
	return r.jwtHandler.MiddlewareFunc()
}
