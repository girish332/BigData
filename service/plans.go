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
	GetAnyObject(c *gin.Context, key string) (interface{}, error)
	GetPlan(c *gin.Context, key string) (models.Plan, error)
	CreatePlan(c *gin.Context, plan models.Plan) error
	DeletePlan(c *gin.Context, objectId string) error
	GetAllPlans(ctx *gin.Context) ([]models.Plan, error)
	PatchPlan(c *gin.Context, key string, plan models.Plan) error
	UpdatePlan(c *gin.Context, objectId string, plan models.Plan) error
}

func (ps *PlansService) GetAnyObject(c *gin.Context, key string) (interface{}, error) {
	value, err := ps.repo.Get(c, key)
	if err != nil {
		log.Printf("Error getting the plan from the redis : %v", err)
		return nil, err
	}

	var plan models.Plan
	err = json.Unmarshal([]byte(value), &plan)
	if err != nil {
		log.Printf("Error unmarshalling the plan from the redis : %v", err)
		return nil, err
	}

	// Check the ObjectType and return the corresponding struct
	switch plan.ObjectType {
	case "membercostshare":
		var pcs models.PlanCostShares
		err = json.Unmarshal([]byte(value), &pcs)
		if err != nil {
			return nil, err
		}
		return pcs, nil
	case "service":
		var ls models.LinkedService
		err = json.Unmarshal([]byte(value), &ls)
		if err != nil {
			return nil, err
		}
		return ls, nil
	case "PlanServiceCostShares":
		var pscs models.PlanServiceCostShares
		err = json.Unmarshal([]byte(value), &pscs)
		if err != nil {
			return nil, err
		}
		return pscs, nil
	case "planservice":
		var lps models.LinkedPlanService
		err = json.Unmarshal([]byte(value), &lps)
		if err != nil {
			return nil, err
		}
		return lps, nil
	default:
		return plan, nil
	}
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

func (ps *PlansService) DeletePlan(c *gin.Context, objectId string) error {
	// Fetch the plan
	plan, err := ps.GetPlan(c, objectId)
	if err != nil {
		log.Printf("Error getting the plan from the redis : %v", err)
		return err
	}

	// Delete the plan
	err = ps.repo.Delete(c, objectId)
	if err != nil {
		log.Printf("Error deleting the plan from the redis : %v", err)
		return err
	}

	// Delete the PlanCostShares
	err = ps.repo.Delete(c, plan.PlanCostShares.ObjectId)
	if err != nil {
		log.Printf("Error deleting the PlanCostShares from the redis : %v", err)
		return err
	}

	// Delete each LinkedPlanService and its related objects
	for _, linkedPlanService := range plan.LinkedPlanServices {
		// Delete the LinkedPlanService
		err = ps.repo.Delete(c, linkedPlanService.ObjectId)
		if err != nil {
			log.Printf("Error deleting the LinkedPlanService from the redis : %v", err)
			//return err
		}

		// Delete the LinkedService
		err = ps.repo.Delete(c, linkedPlanService.LinkedService.ObjectId)
		if err != nil {
			log.Printf("Error deleting the LinkedService from the redis : %v", err)
			return err
		}

		// Delete the PlanServiceCostShares
		err = ps.repo.Delete(c, linkedPlanService.PlanServiceCostShares.ObjectId)
		if err != nil {
			log.Printf("Error deleting the PlanServiceCostShares from the redis : %v", err)
			return err
		}
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

func (ps *PlansService) PatchPlan(c *gin.Context, key string, newPlan models.Plan) error {
	existingPlan, err := ps.GetPlan(c, key)
	if err != nil || existingPlan.ObjectId == "" {
		return err
	}

	// Create a map of new LinkedPlanServices for easy lookup
	newLinkedPlanServices := make(map[string]models.LinkedPlanService)
	for _, newLinkedPlanService := range newPlan.LinkedPlanServices {
		newLinkedPlanServices[newLinkedPlanService.ObjectId] = newLinkedPlanService
	}

	// Update existing LinkedPlanServices if they are in the newLinkedPlanServices map
	for i, existingLinkedPlanService := range existingPlan.LinkedPlanServices {
		if newLinkedPlanService, ok := newLinkedPlanServices[existingLinkedPlanService.ObjectId]; ok {
			existingPlan.LinkedPlanServices[i] = newLinkedPlanService
			delete(newLinkedPlanServices, existingLinkedPlanService.ObjectId)
		}
	}

	// Append any remaining new LinkedPlanServices that were not in the existing plan
	for _, newLinkedPlanService := range newLinkedPlanServices {
		existingPlan.LinkedPlanServices = append(existingPlan.LinkedPlanServices, newLinkedPlanService)
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

func (ps *PlansService) UpdatePlan(c *gin.Context, objectId string, plan models.Plan) error {
	// Delete the existing plan and all its associated objects
	err := ps.DeletePlan(c, objectId)
	if err != nil {
		log.Printf("Failed to delete existing plan with error : %v", err.Error())
		return err
	}

	// Create a new plan with the new request body
	err = ps.CreatePlan(c, plan)
	if err != nil {
		log.Printf("Failed to create new plan with error : %v", err.Error())
		return err
	}

	return nil
}
