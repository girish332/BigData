package handler

import (
	"BigData/models"
	"BigData/service"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Handler interface {
	CreatePlan(c *gin.Context)
	GetPlan(c *gin.Context)
	GetAllPlans(c *gin.Context)
	DeletePlan(c *gin.Context)
	PatchPlan(c *gin.Context)
	UpdatePlan(c *gin.Context)
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

	// Check if a plan with the same objectId already exists
	existingPlan, err := ph.service.GetPlan(c, planRequest.ObjectId)
	if err == nil && existingPlan.ObjectId != "" {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	err = ph.service.CreatePlan(c, planRequest)
	if err != nil {
		log.Printf("Failed to create plan with error : %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//ut.PrettyPrints(planRequest)
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
	clientEtag := strings.TrimSpace(c.GetHeader("If-None-Match"))

	plan, err := ph.service.GetAnyObject(c, objectId)
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
		if err.Error() == "key not found" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Printf("Failed to delete plan with err : %v", err.Error())
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Status(204)
	log.Printf("Plan with objectId : %s deleted successfully", objectId)
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

func (ph *PlansHandler) PatchPlan(c *gin.Context) {
	objectId, ok := c.Params.Get("objectId")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var planRequest models.Plan
	err := c.ShouldBindBodyWith(&planRequest, binding.JSON)
	if err != nil {
		log.Printf("Bad Request with error : %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	existingPlan, err := ph.service.GetPlan(c, objectId)
	if err != nil || existingPlan.ObjectId == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = ph.service.PatchPlan(c, objectId, planRequest)
	if err != nil {
		log.Printf("Failed to update plan with error : %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan updated successfully"})
	return
}

func generateETag(plan interface{}) string {
	h := sha1.New()
	dataBytes, err := json.Marshal(plan)
	if err != nil {
		log.Printf("Error marshalling the plan struct : %v", err)
		return ""
	}
	h.Write(dataBytes)
	sha1Hash := hex.EncodeToString(h.Sum(nil))

	return sha1Hash
}

func (ph *PlansHandler) UpdatePlan(c *gin.Context) {
	var planRequest models.Plan
	err := c.ShouldBindBodyWith(&planRequest, binding.JSON)
	if err != nil {
		log.Printf("Bad Request with error : %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	existingPlan, err := ph.service.GetPlan(c, planRequest.ObjectId)
	if err != nil || existingPlan.ObjectId == "" {
		// If the plan does not exist, create a new one
		err = ph.service.CreatePlan(c, planRequest)
		if err != nil {
			log.Printf("Failed to create plan with error : %v", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Plan created successfully"})
		return
	}

	// If the ETag of the existing plan and the request are the same, return a 304 status
	requestETag := generateETag(planRequest)
	existingETag := generateETag(existingPlan)
	if requestETag == existingETag {
		c.Status(http.StatusNotModified)
		return
	}

	err = ph.service.UpdatePlan(c, planRequest.ObjectId, planRequest)
	if err != nil {
		log.Printf("Failed to update plan with error : %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan updated successfully"})
	return
}
