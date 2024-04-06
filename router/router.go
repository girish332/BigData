package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/girish332/bigdata/database"
	"github.com/girish332/bigdata/elastic"
	"github.com/girish332/bigdata/handler"
	"github.com/girish332/bigdata/middleware"
	"github.com/girish332/bigdata/service"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(gin.Recovery())

	redisRepo := database.NewRedisRepo("localhost:6379", "")
	planService := service.NewPlansService(redisRepo)
	esFactory := elastic.NewElasticFactory()
	planHandler := handler.NewPlansHandler(planService, esFactory)

	v1 := router.Group("/v1", middleware.OAuth2Middleware())
	{
		v1.POST("/plan", planHandler.CreatePlan)
		v1.GET("/plan/:objectId", planHandler.GetPlan)
		v1.DELETE("/plan/:objectId", planHandler.DeletePlan)
		v1.GET("/plans", planHandler.GetAllPlans)
		v1.PATCH("/plan/:objectId", planHandler.PatchPlan)
		v1.PUT("/plan", planHandler.UpdatePlan)
		v1.POST("/search", planHandler.SearchPlans)
	}

	return router
}
