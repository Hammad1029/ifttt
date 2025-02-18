package cron

import "ifttt/manager/domain/api"

type Cron struct {
	Name        string   `json:"name" mapstructure:"name"`
	Description string   `json:"description" mapstructure:"description"`
	CronExpr    string   `json:"cronExpr" mapstructure:"cronExpr"`
	ApiName     string   `json:"apiName" maptructure:"apiName"`
	API         *api.Api `json:"api" mapstructure:"api"`
}
