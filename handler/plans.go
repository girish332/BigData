package handler

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/girish332/bigdata/elastic"
	"github.com/girish332/bigdata/models"
	"github.com/girish332/bigdata/service"
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
	SearchPlans(c *gin.Context)
}

type PlansHandler struct {
	service   *service.PlansService
	esFactory *elastic.Factory
}

func NewPlansHandler(planService *service.PlansService, esFactory *elastic.Factory) *PlansHandler {
	return &PlansHandler{
		service:   planService,
		esFactory: esFactory,
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

	clientEtag := strings.TrimSpace(c.GetHeader("If-None-Match"))
	if clientEtag == "" {
		c.AbortWithStatus(http.StatusPreconditionFailed)
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

	currentEtag := generateETag(planRequest)
	if clientEtag == currentEtag {
		c.Status(http.StatusNotModified)
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

func (h *PlansHandler) SearchPlans(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a match query
	matchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				req.Key: req.Value,
			},
		},
	}
	queryBytes, err := json.Marshal(matchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new search request
	searchReq := esapi.SearchRequest{
		Index: []string{"plans"},
		Body:  bytes.NewReader(queryBytes),
	}

	// Create a new Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // replace with your Elasticsearch address
		},
	}
	client, err := h.esFactory.NewClient(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Perform the search request.
	res, err := searchReq.Do(context.Background(), client.ES)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.String()})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
