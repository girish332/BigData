package router

import (
	"BigData/database"
	"BigData/handler"
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

	router.POST("/v1/plan", planHandler.CreatePlan)
	router.GET("/v1/plan/:objectId", planHandler.GetPlan)
	router.DELETE("/v1/plan/:objectId", planHandler.DeletePlan)
	router.GET("/v1/plans", planHandler.GetAllPlans)

	return router
}
