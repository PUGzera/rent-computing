package machine_rest_handler

import (
	"fmt"
	http_util "rent-computing/internal"
	machine_controller "rent-computing/internal/machine/controller"
	user_controller "rent-computing/internal/user/controller"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type createMachineReq struct {
	Vnc      bool   `json:"vnc"`
	Password string `json:"password"`
}

type machineRestHandler struct {
	router            *gin.RouterGroup
	jwtHandler        *jwt.GinJWTMiddleware
	machineController machine_controller.MachineController
	userController    user_controller.UserController
}

type Options struct {
	Router            *gin.RouterGroup
	JWTHandler        *jwt.GinJWTMiddleware
	UserController    user_controller.UserController
	MachineController machine_controller.MachineController
}

func New(options Options) (*machineRestHandler, error) {
	machineRestHandler := machineRestHandler{
		router:            options.Router,
		jwtHandler:        options.JWTHandler,
		userController:    options.UserController,
		machineController: options.MachineController,
	}

	machineRestHandler.createMachineRoute()
	machineRestHandler.deleteMachineRoute()
	machineRestHandler.startMachineRoute()
	machineRestHandler.stopMachineRoute()
	machineRestHandler.getMachinesRoute()

	return &machineRestHandler, nil
}

func (m *machineRestHandler) createMachineRoute() {
	m.router.
		Use(m.AuthenticationMiddleware()).
		POST("/create", func(ctx *gin.Context) {
			var machineReq createMachineReq
			err := ctx.ShouldBindJSON(&machineReq)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			id, err := http_util.GetUserIdFromCtx(ctx, *m.jwtHandler)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			profile, err := m.userController.Profile(ctx, id)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			machine, err := m.machineController.Create(ctx, id, machineReq.Vnc, profile.Username, machineReq.Password)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(200, machine)
		})
}

func (m *machineRestHandler) getMachinesRoute() {
	m.router.
		Use(m.AuthenticationMiddleware()).
		GET("/", func(ctx *gin.Context) {
			id, err := http_util.GetUserIdFromCtx(ctx, *m.jwtHandler)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			machines, err := m.machineController.GetMachines(ctx, id)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(200, machines)
		})
}

func (m *machineRestHandler) startMachineRoute() {
	m.router.
		Use(m.AuthenticationMiddleware()).
		GET("/start/:id", func(ctx *gin.Context) {
			machineId := ctx.Param("id")

			id, err := http_util.GetUserIdFromCtx(ctx, *m.jwtHandler)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			machine, err := m.machineController.Start(ctx, machineId, id)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(200, machine)
		})
}

func (m *machineRestHandler) stopMachineRoute() {
	m.router.
		Use(m.AuthenticationMiddleware()).
		GET("/stop/:id", func(ctx *gin.Context) {
			machineId := ctx.Param("id")

			id, err := http_util.GetUserIdFromCtx(ctx, *m.jwtHandler)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			err = m.machineController.Stop(ctx, machineId, id)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(200, fmt.Sprintf("successfully stopped machine %s", machineId))
		})
}

func (m *machineRestHandler) deleteMachineRoute() {
	m.router.
		Use(m.AuthenticationMiddleware()).
		DELETE("/delete/:id", func(ctx *gin.Context) {
			machineId := ctx.Param("id")

			id, err := http_util.GetUserIdFromCtx(ctx, *m.jwtHandler)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			err = m.machineController.Delete(ctx, id, machineId)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(200, fmt.Sprintf("successfully deleted machine %s", machineId))
		})
}

func (r *machineRestHandler) AuthenticationMiddleware() gin.HandlerFunc {
	return r.jwtHandler.MiddlewareFunc()
}
