package postgres

import (
	"fmt"
	"ifttt/manager/domain/rule"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"gorm.io/gorm"
)

type PostgresTriggerFlowsRepository struct {
	*PostgresBaseRepository
}

func NewPostgresTriggerFlowsRepository(base *PostgresBaseRepository) *PostgresTriggerFlowsRepository {
	return &PostgresTriggerFlowsRepository{PostgresBaseRepository: base}
}

func (p *PostgresTriggerFlowsRepository) GetAllTriggerFlows() (*[]triggerflow.TriggerFlow, error) {
	var pgTFlows []trigger_flows
	if err := p.client.Preload("Class").Preload("Rules").Find(&pgTFlows).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowByName: could not query trigger flows: %s", err)
	}

	var domainTFlows []triggerflow.TriggerFlow
	for _, tf := range pgTFlows {
		dtf, err := tf.toDomain()
		if err != nil {
			return nil, fmt.Errorf("method *PostgresApiRepository.GetAllTriggerFlows: could not decode trigger flows: %s", err)
		}
		domainTFlows = append(domainTFlows, *dtf)
	}

	return &domainTFlows, nil
}

func (p *PostgresTriggerFlowsRepository) GetTriggerFlowsByNames(names []string) (*[]triggerflow.TriggerFlow, error) {
	var pgTFlows []trigger_flows
	if err := p.client.Find(&pgTFlows, "name in ?", names).Error; err != nil {
		return nil, fmt.Errorf("could not query trigger flows: %s", err)
	}

	var domainTFlows []triggerflow.TriggerFlow
	for _, tf := range pgTFlows {
		dtf, err := tf.toDomain()
		if err != nil {
			return nil, fmt.Errorf("could not decode trigger flows: %s", err)
		}
		domainTFlows = append(domainTFlows, *dtf)
	}

	return &domainTFlows, nil
}

func (p *PostgresTriggerFlowsRepository) GetTriggerFlowByName(name string) (*triggerflow.TriggerFlow, error) {
	var pgTFlow trigger_flows
	if err := p.client.First(&pgTFlow, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowByName: could not query trigger flows: %s", err)
	}

	domainTFlow, err := pgTFlow.toDomain()
	if err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowByName: could not decode trigger flows: %s", err)
	}

	return domainTFlow, nil
}

func (p *PostgresTriggerFlowsRepository) GetTriggerFlowDetailsByName(name string) (*triggerflow.TriggerFlow, error) {
	var pgTFlow trigger_flows
	if err := p.client.
		Preload("Class").Preload("StartRules").Preload("AllRules").
		First(&pgTFlow, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil,
			fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowDetailsByName: could not query trigger flows: %s", err)
	}

	domainTFlow, err := pgTFlow.toDomain()
	if err != nil {
		return nil,
			fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowDetailsByName: could not decode trigger flows: %s", err)
	}

	return domainTFlow, nil
}

func (p *PostgresTriggerFlowsRepository) InsertTriggerFlow(
	tFlow *triggerflow.CreateTriggerFlowRequest, attachRules *[]rule.Rule) error {
	var pgTFlow trigger_flows
	pgTFlow.fromDomain(tFlow, attachRules)
	if err := p.client.Create(&pgTFlow).Error; err != nil {
		return fmt.Errorf("method *PostgresApiRepository.InsertTriggerFlow: could not insert trigger flows: %s", err)
	}
	return nil
}
