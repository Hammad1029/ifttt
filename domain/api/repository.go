package api

import triggerflow "ifttt/manager/domain/trigger_flow"

type Repository interface {
	GetAllApis() (*[]Api, error)
	GetApiByNameOrPath(name string, path string) (*Api, error)
	InsertApi(apiReq *CreateApiRequest, attachTriggers *[]triggerflow.TriggerFlow) error
}
