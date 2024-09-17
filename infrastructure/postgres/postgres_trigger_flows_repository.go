package postgres

import (
	"fmt"
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
	if err := p.client.Preload("Class").Find(&pgTFlows).Error; err != nil {
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

func (p *PostgresTriggerFlowsRepository) GetTriggerFlowsByIds(ids []uint) (*[]triggerflow.TriggerFlow, error) {
	var pgTFlows []trigger_flows
	if err := p.client.Find(&pgTFlows, "id in ?", ids).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowsByIds: could not query trigger flows: %s", err)
	}

	var domainTFlows []triggerflow.TriggerFlow
	for _, tf := range pgTFlows {
		dtf, err := tf.toDomain()
		if err != nil {
			return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowsByIds: could not decode trigger flows: %s", err)
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

func (p *PostgresTriggerFlowsRepository) InsertTriggerFlow(tFlow *triggerflow.CreateTriggerFlowRequest) error {
	var pgTFlow trigger_flows
	pgTFlow.fromDomain(tFlow)
	if err := p.client.Create(&pgTFlow).Error; err != nil {
		return fmt.Errorf("method *PostgresApiRepository.InsertTriggerFlow: could not insert trigger flows: %s", err)
	}
	return nil
}
