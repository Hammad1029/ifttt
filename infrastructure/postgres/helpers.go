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
		Id:          pgRule.ID,
		Name:        pgRule.Name,
		Description: pgRule.Description,
	}

	if err := json.Unmarshal(pgRule.Pre.Bytes, &domainRule.Pre); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling pre: %s", err)
	}

	if err := json.Unmarshal(pgRule.Switch.Bytes, &domainRule.Switch); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling switch: %s", err)
	}

	return &domainRule, nil
}

func (pgRule *rules) fromDomain(domainRule *rule.CreateRuleRequest) error {
	pgRule.Name = domainRule.Name
	pgRule.Description = domainRule.Description

	if preMarshalled, err := json.Marshal(domainRule.Pre); err != nil {
		return fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal conditions: %s", err)
	} else {
		pgRule.Pre = pgtype.JSONB{Bytes: preMarshalled, Status: pgtype.Present}
	}

	if switchMarshalled, err := json.Marshal(domainRule.Switch); err != nil {
		return fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal switch: %s", err)
	} else {
		pgRule.Switch = pgtype.JSONB{Bytes: switchMarshalled, Status: pgtype.Present}
	}

	return nil
}

func (t *trigger_flows) fromDomain(domainTFlow *triggerflow.CreateTriggerFlowRequest) error {
	t.Name = domainTFlow.Name
	t.Description = domainTFlow.Description
	t.ClassId = domainTFlow.Class
	for _, r := range domainTFlow.StartRules {
		t.StartRules = append(t.StartRules, rules{Model: gorm.Model{ID: r}})
	}
	for _, r := range domainTFlow.AllRules {
		t.AllRules = append(t.AllRules, rules{Model: gorm.Model{ID: r}})
	}
	if bfMarshalled, err := json.Marshal(domainTFlow.BranchFlows); err != nil {
		return fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal branchFlow: %s", err)
	} else {
		t.BranchFlow = pgtype.JSONB{Bytes: bfMarshalled, Status: pgtype.Present}
	}
	return nil
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

	if err := json.Unmarshal(t.BranchFlow.Bytes, &domanTFlow.BranchFlows); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling branchFlows: %s", err)
	}

	return &domanTFlow, nil
}

func (a *apis) fromDomain(domainApi *api.CreateApiRequest) error {
	a.Name = domainApi.Name
	a.Path = domainApi.Path
	a.Method = domainApi.Method
	a.Description = domainApi.Description

	if reqMarshalled, err := json.Marshal(domainApi.Request); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal request: %s", err)
	} else {
		a.Request = pgtype.JSONB{Bytes: reqMarshalled, Status: pgtype.Present}
	}

	if preConfigMarshalled, err := json.Marshal(domainApi.PreConfig); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal pre config: %s", err)
	} else {
		a.PreConfig = pgtype.JSONB{Bytes: preConfigMarshalled, Status: pgtype.Present}
	}

	for _, dtf := range domainApi.TriggerFlows {
		a.TriggerFlowRef = append(a.TriggerFlowRef, trigger_flows{Model: gorm.Model{ID: dtf.Trigger}})
	}

	if tConditionsMarshalled, err := json.Marshal(domainApi.TriggerFlows); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal trigger conditions: %s", err)
	} else {
		a.TriggerFlows = pgtype.JSONB{Bytes: tConditionsMarshalled, Status: pgtype.Present}
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

	var tConditions []api.TriggerConditionRequest
	if err := json.Unmarshal(a.TriggerFlows.Bytes, &tConditions); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.ToDomain: could not cast pgApi: %s", err)
	}

	triggerFlowMap := make(map[uint]trigger_flows)
	for _, tFlow := range a.TriggerFlowRef {
		triggerFlowMap[tFlow.ID] = tFlow
	}

	for _, tc := range tConditions {
		tcModel, ok := triggerFlowMap[tc.Trigger]
		if !ok {
			return nil,
				fmt.Errorf("method *PostgresAPIRepository.ToDomain: trigger flow not found from conditions")
		}
		domainTFlow, err := tcModel.toDomain()
		if err != nil {
			return nil, fmt.Errorf("method *PostgresAPIRepository.ToDomain: %s", err)
		}
		*domainApi.TriggerFlows = append(*domainApi.TriggerFlows, api.TriggerCondition{If: tc.If, Trigger: *domainTFlow})
	}

	return &domainApi, nil
}
