package routers

import (
	"crazyfarmbackend/config/di"
	"github.com/gin-gonic/gin"
)

func pong(c *gin.Context) {
	c.String(200, "pong")
}

func SetupMainGroup(router *gin.Engine, init *di.Initialization) {
	api := router.Group("/api/v1")
	{
		api.GET("/ping", pong)
		api.GET("/user/auth", init.UserController.AuthUser)
		api.GET("/user/me", init.MiddlewareService.AuthMiddleware(), init.UserController.GetMe)
		api.GET("/user/upgrade", init.MiddlewareService.AuthMiddleware(), init.UserController.GetMyUpgrades)
		api.GET("/user/fields", init.MiddlewareService.AuthMiddleware(), init.UserController.GetMyFields)
		api.GET("/user/referrals", init.MiddlewareService.AuthMiddleware(), init.UserController.GetMyReferrals)
		api.GET("/inventory/all", init.MiddlewareService.AuthMiddleware(), init.InventoryController.GetInventoryItems)
		api.GET("/tasks/all", init.MiddlewareService.AuthMiddleware(), init.TaskController.GetAllTasks)
		api.GET("/tasks/check", init.MiddlewareService.AuthMiddleware(), init.TaskController.Check)
		api.GET("/tasks/claim", init.MiddlewareService.AuthMiddleware(), init.TaskController.Claim)
	}
}
