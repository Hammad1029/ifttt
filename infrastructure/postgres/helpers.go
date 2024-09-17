package postgres

import (
	"encoding/json"
	"fmt"
	"ifttt/manager/domain/api"
	"ifttt/manager/domain/rule"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

func (pgRule *rules) toDomain() (*rule.Rule, error) {
	domainRule := rule.Rule{
		Name:        pgRule.Name,
		Description: pgRule.Description,
	}

	if err := json.Unmarshal(pgRule.Conditions.Bytes, &domainRule.Conditions); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling conditions: %s", err)
	}

	if err := json.Unmarshal(pgRule.Then.Bytes, &domainRule.Then); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling then: %s", err)
	}

	if err := json.Unmarshal(pgRule.Else.Bytes, &domainRule.Else); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling else: %s", err)
	}

	return &domainRule, nil
}

func (pgRule *rules) fromDomain(domainRule *rule.CreateRuleRequest) error {
	pgRule.Name = domainRule.Name
	pgRule.Description = domainRule.Description

	if conditionsMarshalled, err := json.Marshal(domainRule.Conditions); err != nil {
		return fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal conditions: %s", err)
	} else {
		pgRule.Conditions = pgtype.JSONB{Bytes: conditionsMarshalled, Status: pgtype.Present}
	}

	if thenMarshalled, err := json.Marshal(domainRule.Then); err != nil {
		return fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal then: %s", err)
	} else {
		pgRule.Then = pgtype.JSONB{Bytes: thenMarshalled, Status: pgtype.Present}
	}

	if elseMarshalled, err := json.Marshal(domainRule.Else); err != nil {
		return fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal else: %s", err)
	} else {
		pgRule.Else = pgtype.JSONB{Bytes: elseMarshalled, Status: pgtype.Present}
	}

	return nil
}

func (t *trigger_flows) fromDomain(domainTFlow *triggerflow.CreateTriggerFlowRequest) {
	t.Name = domainTFlow.Name
	t.Description = domainTFlow.Description
	t.ClassId = domainTFlow.Class
	for _, r := range domainTFlow.StartRules {
		t.StartRules = append(t.StartRules, rules{Model: gorm.Model{ID: r}})
	}
	for _, r := range domainTFlow.AllRules {
		t.AllRules = append(t.AllRules, rules{Model: gorm.Model{ID: r}})
	}
}

func (t *trigger_flows) toDomain() (*triggerflow.TriggerFlow, error) {
	domanTFlow := triggerflow.TriggerFlow{
		Name:        t.Name,
		Description: t.Description,
		Class:       triggerflow.Class{Name: t.Class.Name},
	}
	for _, r := range t.StartRules {
		dRule, err := r.toDomain()
		if err != nil {
			return nil,
				fmt.Errorf("method *PostgresTriggerFlowsRepository.ToDomain: could not convert to domain rule")
		}
		domanTFlow.StartRules = append(domanTFlow.StartRules, *dRule)
	}
	for _, r := range t.AllRules {
		dRule, err := r.toDomain()
		if err != nil {
			return nil,
				fmt.Errorf("method *PostgresTriggerFlowsRepository.ToDomain: could not convert to domain rule")
		}
		domanTFlow.AllRules = append(domanTFlow.AllRules, *dRule)
	}

	return &domanTFlow, nil
}

func (a *apis) fromDomain(domainApi *api.CreateApiRequest) error {
	a.Name = domainApi.Name
	a.Path = domainApi.Path
	a.Method = domainApi.Method
	a.Description = domainApi.Description

	if reqMarshalled, err := json.Marshal(domainApi.Request); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal request")
	} else {
		a.Request = pgtype.JSONB{Bytes: reqMarshalled, Status: pgtype.Present}
	}

	if preConfigMarshalled, err := json.Marshal(domainApi.PreConfig); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal pre config")
	} else {
		a.PreConfig = pgtype.JSONB{Bytes: preConfigMarshalled, Status: pgtype.Present}
	}

	for _, dtf := range domainApi.TriggerFlows {
		a.Triggerflows = append(a.Triggerflows, trigger_flows{Model: gorm.Model{ID: dtf}})
	}

	return nil
}

func (a *apis) toDomain() (*api.Api, error) {
	domainApi := api.Api{
		Name:        a.Name,
		Path:        a.Path,
		Method:      a.Method,
		Description: a.Description,
	}

	if err := json.Unmarshal(a.Request.Bytes, &domainApi.Request); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.ToDomain: could not cast pgApi: %s", err)
	}

	if err := json.Unmarshal(a.PreConfig.Bytes, &domainApi.PreConfig); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.ToDomain: could not cast pgApi: %s", err)
	}

	return &domainApi, nil
}
