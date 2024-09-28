package api

import (
	"crazyfarmbackend/config/di"
	"crazyfarmbackend/src/api/routers"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Init(init *di.Initialization) *gin.Engine {
	router := gin.New()
	router.Use(init.MiddlewareService.RequestIdMiddleware())
	router.Use(init.MiddlewareService.CorsMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.BestSpeed))
	{
		routers.SetupMainGroup(router, init)
	}
	return router
}
