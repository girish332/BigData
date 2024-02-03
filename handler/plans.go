package handler

import (
	"BigData/models"
	"BigData/service"
	ut "BigData/utils"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Handler interface {
	CreatePlan(c *gin.Context)
	GetPlan(c *gin.Context)
	GetAllPlans(c *gin.Context)
	DeletePlan(c *gin.Context)
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
	eTag := generateETag(planRequest)
	c.Header("ETag", eTag)
	c.JSON(http.StatusCreated, gin.H{"message": "Plan created successfully"})
	return
}

func (ph *PlansHandler) GetPlan(c *gin.Context) {
	objectId, ok := c.Params.Get("objectId")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	clientEtag := c.GetHeader("If-None-Match")

	plan, err := ph.service.GetPlan(c, objectId)
	if err != nil {
		log.Printf("Failed to fetch plan with err : %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	currentEtag := generateETag(plan)
	if clientEtag == currentEtag {
		c.Status(http.StatusNotModified)
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
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Status(204)
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

func generateETag(plan models.Plan) string {
	h := sha1.New()
	dataBytes, err := json.Marshal(plan)
	if err != nil {
		log.Printf("Error marshalling the plan struct : %v", err)
		return ""
	}
	h.Write(dataBytes)
	sha1Hash := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("\"%s\"", sha1Hash)
}
