package user_rest_handler

import (
	http_util "rent-computing/internal"
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
	Id        string    `json:"id"`
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
	JWTHandler     *jwt.GinJWTMiddleware
	UserController user_controller.UserController
}

func New(options Options) (*userRestHandler, error) {
	userRestHandler := userRestHandler{
		router:         options.Router,
		jwtHandler:     options.JWTHandler,
		userController: options.UserController,
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
			Id:        user.Id,
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
	r.router.
		Use(r.AuthenticationMiddleware()).
		GET("/profile", func(ctx *gin.Context) {
			user, err := r.GetAuthenticatedUser(ctx)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, gin.H{"user": user})
		})
}

func (r *userRestHandler) GetAuthenticatedUser(ctx *gin.Context) (*authenticatedUser, error) {
	id, err := http_util.GetUserIdFromCtx(ctx, *r.jwtHandler)
	if err != nil {
		return nil, err
	}

	user, err := r.userController.Profile(ctx, id)
	if err != nil {
		return nil, err
	}

	return &authenticatedUser{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (r *userRestHandler) AuthenticationMiddleware() gin.HandlerFunc {
	return r.jwtHandler.MiddlewareFunc()
}
