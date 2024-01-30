package handler

import (
	"BigData/models"
	"BigData/service"
	ut "BigData/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Handler interface {
	CreatePlan(c *gin.Context)
	GetPlan(c *gin.Context)
}

type PlansHandler struct {
	service *service.PlansService
}

func NewPlansHandler(planService *service.PlansService) *PlansHandler {
	return &PlansHandler{
		service: planService,
	}
}

func (ph *PlansHandler) CreatePlan(c *gin.Context) {
	var planRequest models.Plan
	err := c.ShouldBindBodyWith(&planRequest, binding.JSON)
	if err != nil {
		log.Printf("Bad Request with error : %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = ph.service.CreatePlan(c, planRequest)
	if err != nil {
		log.Printf("Failed to create plan with error : %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ut.PrettyPrints(planRequest)
	c.JSON(http.StatusCreated, gin.H{"message": "Plan created successfully"})
	return
}

func (ph *PlansHandler) GetPlan(c *gin.Context) {
	objectId, ok := c.Params.Get("objectId")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	plan, err := ph.service.GetPlan(c, objectId)
	if err != nil {
		log.Printf("Failed to fetch plan with err : %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, plan)
	return
}

func (ph *PlansHandler) DeletePlan(c *gin.Context) {
	objectId, ok := c.Params.Get("objectId")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := ph.service.DeletePlan(c, objectId)
	if err != nil {
		log.Printf("Failed to delete plan with err : %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted successfully"})
	return
}

func (ph *PlansHandler) GetAllPlans(c *gin.Context) {
	plans, err := ph.service.GetAllPlans(c)
	if err != nil {
		log.Printf("Failed to fetch all plans with err : %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, plans)
	return
}
