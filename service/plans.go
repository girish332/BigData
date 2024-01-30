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

	return ps.repo.Set(c, objectId, string(value))
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
		plan, err := ps.GetPlan(ctx, key)
		if err != nil {
			log.Printf("Error fetching the plan from the redis : %v", err)
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, nil
}
