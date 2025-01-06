package cron

import triggerflow "ifttt/manager/domain/trigger_flow"

type Repository interface {
	GetAllCrons() (*[]Cron, error)
	GetCronByName(name string) (*Cron, error)
	InsertCron(req *CreateCronRequest, attachTriggers *[]triggerflow.TriggerFlow) error
}
