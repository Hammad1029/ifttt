package postgres

import (
	"encoding/json"
	"fmt"
	"ifttt/manager/domain/api"
	"ifttt/manager/domain/cron"
	eventprofiles "ifttt/manager/domain/event_profiles"
	"ifttt/manager/domain/orm_schema"
	requestvalidator "ifttt/manager/domain/request_validator"
	"ifttt/manager/domain/resolvable"
	"ifttt/manager/domain/rule"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"github.com/go-viper/mapstructure/v2"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

func (pgRule *rules) toDomain() (*rule.Rule, error) {
	domainRule := rule.Rule{
		ID:          pgRule.ID,
		Name:        pgRule.Name,
		Description: pgRule.Description,
	}

	if err := json.Unmarshal(pgRule.Pre.Bytes, &domainRule.Pre); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(pgRule.Switch.Bytes, &domainRule.Switch); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(pgRule.Finally.Bytes, &domainRule.Finally); err != nil {
		return nil, err
	}

	return &domainRule, nil
}

func (pgRule *rules) fromDomain(domainRule *rule.CreateRuleRequest) error {
	pgRule.Name = domainRule.Name
	pgRule.Description = domainRule.Description

	if preMarshalled, err := json.Marshal(domainRule.Pre); err != nil {
		return fmt.Errorf("could not marshal Pre: %s", err)
	} else {
		pgRule.Pre = pgtype.JSONB{Bytes: preMarshalled, Status: pgtype.Present}
	}

	if switchMarshalled, err := json.Marshal(domainRule.Switch); err != nil {
		return fmt.Errorf("could not marshal switch: %s", err)
	} else {
		pgRule.Switch = pgtype.JSONB{Bytes: switchMarshalled, Status: pgtype.Present}
	}

	if finallyMarshalled, err := json.Marshal(domainRule.Finally); err != nil {
		return fmt.Errorf("could not marshal finally: %s", err)
	} else {
		pgRule.Finally = pgtype.JSONB{Bytes: finallyMarshalled, Status: pgtype.Present}
	}

	return nil
}

func (t *trigger_flows) fromDomain(domainTFlow *triggerflow.CreateTriggerFlowRequest, attachRules *[]rule.Rule) error {
	t.Name = domainTFlow.Name
	t.Description = domainTFlow.Description
	t.StartState = domainTFlow.StartState
	for _, r := range *attachRules {
		t.Rules = append(t.Rules, rules{Model: gorm.Model{ID: r.ID}})
	}
	if bfMarshalled, err := json.Marshal(domainTFlow.BranchFlows); err != nil {
		return fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal branchFlow: %s", err)
	} else {
		t.BranchFlows = pgtype.JSONB{Bytes: bfMarshalled, Status: pgtype.Present}
	}
	return nil
}

func (t *trigger_flows) toDomain() (*triggerflow.TriggerFlow, error) {
	domanTFlow := triggerflow.TriggerFlow{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		StartState:  t.StartState,
		Rules:       map[uint]*rule.Rule{},
		BranchFlows: map[uint]*triggerflow.BranchFlow{},
	}
	for _, r := range t.Rules {
		dRule, err := r.toDomain()
		if err != nil {
			return nil,
				fmt.Errorf("method *PostgresTriggerFlowsRepository.ToDomain: could not convert to domain rule")
		}
		domanTFlow.Rules[r.ID] = dRule
	}

	if err := json.Unmarshal(t.BranchFlows.Bytes, &domanTFlow.BranchFlows); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling branchFlows: %s", err)
	}

	return &domanTFlow, nil
}

func (a *apis) fromDomain(domainApi *api.CreateApiRequest, attachTriggers *[]triggerflow.TriggerFlow) error {
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

	for _, dtf := range *attachTriggers {
		a.Triggers = append(a.Triggers, trigger_flows{Model: gorm.Model{ID: dtf.ID}})
	}

	if tConditionsMarshalled, err := json.Marshal(domainApi.Triggers); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal trigger conditions: %s", err)
	} else {
		a.TriggerFlows = pgtype.JSONB{Bytes: tConditionsMarshalled, Status: pgtype.Present}
	}

	return nil
}

func (a *apis) toDomain() (*api.Api, error) {
	domainApi := api.Api{
		ID:          a.ID,
		Name:        a.Name,
		Path:        a.Path,
		Method:      a.Method,
		Description: a.Description,
		Request:     map[string]requestvalidator.RequestParameter{},
		PreConfig:   map[string]resolvable.Resolvable{},
		Triggers:    &[]triggerflow.TriggerCondition{},
	}

	if err := json.Unmarshal(a.Request.Bytes, &domainApi.Request); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(a.PreConfig.Bytes, &domainApi.PreConfig); err != nil {
		return nil, err
	}

	var tConditions []triggerflow.TriggerConditionRequest
	if err := json.Unmarshal(a.TriggerFlows.Bytes, &tConditions); err != nil {
		return nil, err
	}

	triggerFlowMap := make(map[string]trigger_flows)
	for _, tFlow := range a.Triggers {
		triggerFlowMap[tFlow.Name] = tFlow
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
		*domainApi.Triggers = append(*domainApi.Triggers,
			triggerflow.TriggerCondition{If: tc.If, Trigger: *domainTFlow})
	}

	return &domainApi, nil
}

func (c *crons) fromDomain(dCron *cron.CreateCronRequest, attachTriggers *[]triggerflow.TriggerFlow) error {
	c.Name = dCron.Name
	c.Description = dCron.Description
	c.Cron = dCron.Cron

	if preConfigMarshalled, err := json.Marshal(dCron.PreConfig); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal pre config: %s", err)
	} else {
		c.PreConfig = pgtype.JSONB{Bytes: preConfigMarshalled, Status: pgtype.Present}
	}

	for _, dtf := range *attachTriggers {
		c.TriggerFlowRef = append(c.TriggerFlowRef, trigger_flows{Model: gorm.Model{ID: dtf.ID}})
	}

	if tConditionsMarshalled, err := json.Marshal(dCron.Triggers); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal trigger conditions: %s", err)
	} else {
		c.TriggerFlows = pgtype.JSONB{Bytes: tConditionsMarshalled, Status: pgtype.Present}
	}

	return nil
}

func (c *crons) toDomain() (*cron.Cron, error) {
	dCron := cron.Cron{
		ID:           c.ID,
		Name:         c.Name,
		Description:  c.Name,
		Cron:         c.Cron,
		TriggerFlows: &[]triggerflow.TriggerCondition{},
	}

	if err := json.Unmarshal(c.PreConfig.Bytes, &dCron.PreConfig); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.ToDomain: could not cast pgApi: %s", err)
	}

	var tConditions []triggerflow.TriggerConditionRequest
	if err := json.Unmarshal(c.TriggerFlows.Bytes, &tConditions); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.ToDomain: could not cast pgApi: %s", err)
	}

	triggerFlowMap := make(map[string]trigger_flows)
	for _, tFlow := range c.TriggerFlowRef {
		triggerFlowMap[tFlow.Name] = tFlow
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
		*dCron.TriggerFlows = append(*dCron.TriggerFlows, triggerflow.TriggerCondition{If: tc.If, Trigger: *domainTFlow})
	}

	return &dCron, nil
}

func (o *orm_model) fromDomain(dModel *orm_schema.Model) error {
	return mapstructure.Decode(dModel, o)
}

func (o *orm_association) fromDomain(dAssociation *orm_schema.ModelAssociation) error {
	return mapstructure.Decode(dAssociation, o)
}

func (o *orm_model) toDomain() (*orm_schema.Model, error) {
	var domain orm_schema.Model
	if err := mapstructure.Decode(o, &domain); err != nil {
		return nil, err
	}
	domain.ID = o.Model.ID
	return &domain, nil
}

func (o *orm_association) toDomain() (*orm_schema.ModelAssociation, error) {
	var domain orm_schema.ModelAssociation
	if err := mapstructure.Decode(o, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}

func (p *event_profile) fromDomain(dProfile *eventprofiles.Profile) error {
	p.Trigger = dProfile.Trigger
	p.ResponseHTTPStatus = dProfile.ResponseHTTPStatus
	p.Internal = dProfile.Internal
	p.UseBody = dProfile.UseBody
	if marshalled, err := json.Marshal(dProfile.ResponseBody); err != nil {
		return err
	} else {
		p.ResponseBody = pgtype.JSONB{Bytes: marshalled, Status: pgtype.Present}
	}
	return nil
}

func (p *event_profile) toDomain() (*eventprofiles.Profile, error) {
	dProfile := eventprofiles.Profile{
		ID:                 p.ID,
		Trigger:            p.Trigger,
		ResponseHTTPStatus: p.ResponseHTTPStatus,
		Internal:           p.Internal,
		UseBody:            p.UseBody,
	}
	if err := json.Unmarshal(p.ResponseBody.Bytes, &dProfile.ResponseBody); err != nil {
		return nil, err
	}
	if p.MappedProfiles != nil {
		mappedProfiles := make([]eventprofiles.Profile, 0, len(*p.MappedProfiles))
		for _, mp := range *p.MappedProfiles {
			if dP, err := mp.toDomain(); err != nil {
				return nil, err
			} else {
				mappedProfiles = append(mappedProfiles, *dP)
			}
		}
		dProfile.MappedProfiles = &mappedProfiles
	}
	return &dProfile, nil
}
