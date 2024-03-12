package router

import (
	"BigData/database"
	"BigData/handler"
	"BigData/middleware"
	"BigData/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(gin.Recovery())

	redisRepo := database.NewRedisRepo("localhost:6379", "")
	planService := service.NewPlansService(redisRepo)
	planHandler := handler.NewPlansHandler(planService)

	v1 := router.Group("/v1", middleware.OAuth2Middleware())
	{
		v1.POST("/plan", planHandler.CreatePlan)
		v1.GET("/plan/:objectId", planHandler.GetPlan)
		v1.DELETE("/plan/:objectId", planHandler.DeletePlan)
		v1.GET("/plans", planHandler.GetAllPlans)
		v1.PATCH("/plan/:objectId", planHandler.PatchPlan)
	}

	return router
}
