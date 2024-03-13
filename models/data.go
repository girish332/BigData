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

func (plan *Plan) UpdatePlan(updatedPlan Plan) {
	plan.PlanCostShares = updatedPlan.PlanCostShares
	plan.LinkedPlanServices = updatedPlan.LinkedPlanServices
	plan.PlanType = updatedPlan.PlanType
	plan.CreationDate = updatedPlan.CreationDate
	plan.Org = updatedPlan.Org
}

func (l *LinkedPlanService) UpdateLinkedPlanService(updatedLinkedPlanService LinkedPlanService) {
	l.LinkedService = updatedLinkedPlanService.LinkedService
	l.PlanServiceCostShares = updatedLinkedPlanService.PlanServiceCostShares
	l.ObjectId = updatedLinkedPlanService.ObjectId
	l.ObjectType = updatedLinkedPlanService.ObjectType
	l.Org = updatedLinkedPlanService.Org
}

func (p *PlanCostShares) UpdatePlanCostShares(updatedPlanCostShares PlanCostShares) {
	p.Deductible = updatedPlanCostShares.Deductible
	p.Copay = updatedPlanCostShares.Copay
	p.ObjectId = updatedPlanCostShares.ObjectId
	p.ObjectType = updatedPlanCostShares.ObjectType
	p.Org = updatedPlanCostShares.Org
}

func (p *LinkedService) UpdateLinkedService(updatedLinkedService LinkedService) {
	p.ObjectId = updatedLinkedService.ObjectId
	p.ObjectType = updatedLinkedService.ObjectType
	p.Name = updatedLinkedService.Name
	p.Org = updatedLinkedService.Org
}

func (p *PlanServiceCostShares) UpdatePlanServiceCostShares(updatedPlanServiceCostShares PlanServiceCostShares) {
	p.Deductible = updatedPlanServiceCostShares.Deductible
	p.Copay = updatedPlanServiceCostShares.Copay
	p.ObjectId = updatedPlanServiceCostShares.ObjectId
	p.ObjectType = updatedPlanServiceCostShares.ObjectType
	p.Org = updatedPlanServiceCostShares.Org
}
