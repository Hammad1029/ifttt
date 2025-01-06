package triggerflow

import "ifttt/manager/domain/rule"

type Repository interface {
	GetAllTriggerFlows() (*[]TriggerFlow, error)
	GetTriggerFlowsByNames(names []string) (*[]TriggerFlow, error)
	GetTriggerFlowByName(name string) (*TriggerFlow, error)
	GetTriggerFlowDetailsByName(name string) (*TriggerFlow, error)
	InsertTriggerFlow(tFlow *CreateTriggerFlowRequest, attachRules *[]rule.Rule) error
}
