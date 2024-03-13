package service

import (
	"BigData/models"
	"BigData/repository"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PlansService struct {
	repo repository.RedisRepo
}

func NewPlansService(repo repository.RedisRepo) *PlansService {
	return &PlansService{
		repo: repo,
	}
}

type Service interface {
	GetPlan(c *gin.Context, key string) (models.Plan, error)
	CreatePlan(c *gin.Context, plan models.Plan) error
	DeletePlan(c *gin.Context, key string) error
	GetAllPlans(ctx *gin.Context) ([]models.Plan, error)
	PatchPlan(c *gin.Context, key string, plan models.Plan) error
}

func (ps *PlansService) GetPlan(c *gin.Context, key string) (models.Plan, error) {
	value, err := ps.repo.Get(c, key)
	if err != nil {
		log.Printf("Error getting the plan from the redis : %v", err)
		return models.Plan{}, err
	}

	var plan models.Plan
	err = json.Unmarshal([]byte(value), &plan)
	if err != nil {
		log.Printf("Error unmarshalling the plan from the redis : %v", err)
		return models.Plan{}, err
	}

	return plan, nil
}

func (ps *PlansService) CreatePlan(c *gin.Context, plan models.Plan) error {
	// Add Code to marshal the struct into a string and set it in the redis
	objectId := plan.ObjectId

	// Marshal the struct into a string
	value, err := json.Marshal(plan)
	if err != nil {
		log.Errorf("Error marshalling the plan struct : %v", err)
		return err
	}
	err = ps.repo.Set(c, objectId, string(value))
	if err != nil {
		log.Printf("Error setting the plan in the redis : %v", err)
		return err
	}

	pValue, err := json.Marshal(plan.PlanCostShares)
	if err != nil {
		log.Errorf("Error marshalling the plan struct : %v", err)
		return err
	}
	err = ps.repo.Set(c, plan.PlanCostShares.ObjectId, string(pValue))
	if err != nil {
		log.Printf("Error setting the plan in the redis : %v", err)
		return err
	}

	// Iterate over the linkedPlanServices array and set each object in Redis
	for _, linkedPlanService := range plan.LinkedPlanServices {
		// Marshal and store the entire linkedPlanService object in Redis
		linkedPlanServiceValue, err := json.Marshal(linkedPlanService)
		if err != nil {
			log.Errorf("Error marshalling the linkedPlanService struct : %v", err)
			return err
		}

		err = ps.repo.Set(c, linkedPlanService.ObjectId, string(linkedPlanServiceValue))
		if err != nil {
			return err
		}

		// Marshal and store the linkedService object in Redis
		linkedServiceValue, err := json.Marshal(linkedPlanService.LinkedService)
		if err != nil {
			log.Errorf("Error marshalling the linkedService struct : %v", err)
			return err
		}

		err = ps.repo.Set(c, linkedPlanService.LinkedService.ObjectId, string(linkedServiceValue))
		if err != nil {
			return err
		}

		// Marshal and store the planserviceCostShares object in Redis
		planserviceCostSharesValue, err := json.Marshal(linkedPlanService.PlanServiceCostShares)
		if err != nil {
			log.Errorf("Error marshalling the planserviceCostShares struct : %v", err)
			return err
		}

		err = ps.repo.Set(c, linkedPlanService.PlanServiceCostShares.ObjectId, string(planserviceCostSharesValue))
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *PlansService) DeletePlan(c *gin.Context, key string) error {
	err := ps.repo.Delete(c, key)
	if err != nil {
		log.Printf("Error deleting the plan from the redis : %v", err)
		return err
	}

	return nil
}

func (ps *PlansService) GetAllPlans(ctx *gin.Context) ([]models.Plan, error) {
	plans := make([]models.Plan, 0)
	keys, err := ps.repo.Keys(ctx, "*")
	if err != nil {
		log.Printf("Error fetching all the keys from the redis : %v", err)
		return nil, err
	}
	for _, key := range keys {
		value, err := ps.repo.Get(ctx, key)
		if err != nil {
			log.Printf("Error fetching the value from the redis : %v", err)
			return nil, err
		}

		var plan models.Plan
		err = json.Unmarshal([]byte(value), &plan)
		if err != nil {
			log.Printf("Error unmarshalling the plan from the redis : %v", err)
			continue
		}

		if plan.ObjectId != "" && plan.ObjectType != "" && plan.PlanCostShares.ObjectId != "" {
			plans = append(plans, plan)
		}
	}

	return plans, nil
}

func (ps *PlansService) PatchPlan(c *gin.Context, key string, plan models.Plan) error {
	existingPlan, err := ps.GetPlan(c, key)
	if err != nil || existingPlan.ObjectId == "" {
		return err
	}

	// If there are new linkedPlanServices, check their objectId and append if they are not in the existing ones
	for _, newLinkedPlanService := range plan.LinkedPlanServices {
		exists := false
		for _, existingLinkedPlanService := range existingPlan.LinkedPlanServices {
			if newLinkedPlanService.ObjectId == existingLinkedPlanService.ObjectId {
				exists = true
				break
			}
		}
		if !exists {
			existingPlan.LinkedPlanServices = append(existingPlan.LinkedPlanServices, newLinkedPlanService)
		}
	}

	// Marshal the updated plan into a string
	value, err := json.Marshal(existingPlan)
	if err != nil {
		return err
	}

	// Update the plan in the database
	err = ps.repo.Set(c, key, string(value))
	if err != nil {
		return err
	}

	return nil
}
