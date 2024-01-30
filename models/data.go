package models

type PlanCostShares struct {
	Deductible int    `json:"deductible" binding:"required"`
	Copay      int    `json:"copay" binding:"required"`
	ObjectId   string `json:"objectId" binding:"required"`
	ObjectType string `json:"objectType" binding:"required"`
	Org        string `json:"_org" binding:"required"`
}

type LinkedService struct {
	ObjectId   string `json:"objectId" binding:"required"`
	ObjectType string `json:"objectType" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Org        string `json:"_org" binding:"required"`
}

type PlanServiceCostShares struct {
	Deductible int    `json:"deductible" binding:"required"`
	Copay      int    `json:"copay" binding:"required"`
	ObjectId   string `json:"objectId" binding:"required"`
	ObjectType string `json:"objectType" binding:"required"`
	Org        string `json:"_org" binding:"required"`
}

type LinkedPlanService struct {
	LinkedService         LinkedService         `json:"linkedService" binding:"required"`
	PlanServiceCostShares PlanServiceCostShares `json:"planserviceCostShares" binding:"required"`
	ObjectId              string                `json:"objectId" binding:"required"`
	ObjectType            string                `json:"objectType" binding:"required"`
	Org                   string                `json:"_org" binding:"required"`
}

type Plan struct {
	PlanCostShares     PlanCostShares      `json:"planCostShares" binding:"required"`
	LinkedPlanServices []LinkedPlanService `json:"linkedPlanServices" binding:"required"`
	ObjectId           string              `json:"objectId" binding:"required"`
	ObjectType         string              `json:"objectType" binding:"required"`
	PlanType           string              `json:"planType" binding:"required"`
	CreationDate       string              `json:"creationDate" binding:"required"`
	Org                string              `json:"_org" binding:"required"`
}
