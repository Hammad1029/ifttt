package triggerflow

type Repository interface {
	GetAllTriggerFlows() (*[]TriggerFlow, error)
	GetTriggerFlowsByIds(ids []uint) (*[]TriggerFlow, error)
	GetTriggerFlowByName(name string) (*TriggerFlow, error)
	GetTriggerFlowDetailsByName(name string) (*TriggerFlow, error)
	InsertTriggerFlow(tFlow *CreateTriggerFlowRequest) error
}
