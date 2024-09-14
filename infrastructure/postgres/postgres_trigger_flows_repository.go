package postgres

import (
	"fmt"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"github.com/go-viper/mapstructure/v2"
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
	if err := p.client.Find(&pgTFlows); err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowByName: could not query trigger flows: %s", err)
	}

	var domainTFlows []triggerflow.TriggerFlow
	if err := mapstructure.Decode(pgTFlows, &domainTFlows); err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowByName: could not decode trigger flows: %s", err)
	}

	return &domainTFlows, nil
}

func (p *PostgresTriggerFlowsRepository) GetTriggerFlowsByIds(ids []uint) (*[]triggerflow.TriggerFlow, error) {
	var pgTFlows []trigger_flows
	if err := p.client.Find(&pgTFlows, "id in ?", ids).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowsByIds: could not query trigger flows: %s", err)
	}

	var domainTFlows []triggerflow.TriggerFlow
	if err := mapstructure.Decode(pgTFlows, &domainTFlows); err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowsByIds: could not decode trigger flows: %s", err)
	}

	return &domainTFlows, nil
}

func (p *PostgresTriggerFlowsRepository) GetTriggerFlowByName(name string) (*triggerflow.TriggerFlow, error) {
	var pgTFlow trigger_flows
	if err := p.client.First(&pgTFlow, "name = ?", name).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowByName: could not query trigger flows: %s", err)
	}

	var domainTFlow triggerflow.TriggerFlow
	if err := mapstructure.Decode(pgTFlow, &domainTFlow); err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowByName: could not decode trigger flows: %s", err)
	}

	return &domainTFlow, nil
}

func (p *PostgresTriggerFlowsRepository) GetTriggerFlowDetailsByName(name string) (*triggerflow.TriggerFlow, error) {
	var pgTFlow trigger_flows
	if err := p.client.
		Preload("StartRules").Preload("AllRules").First(&pgTFlow, "name = ?", name).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowDetailsByName: could not query trigger flows: %s", err)
	}

	var domainTFlow triggerflow.TriggerFlow
	if err := mapstructure.Decode(pgTFlow, &domainTFlow); err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetTriggerFlowDetailsByName: could not decode trigger flows: %s", err)
	}

	return &domainTFlow, nil
}

func (p *PostgresTriggerFlowsRepository) InsertTriggerFlow(tFlow *triggerflow.CreateTriggerFlowRequest) error {
	newFlow := trigger_flows{
		Name:        tFlow.Name,
		Description: tFlow.Description,
		ClassId:     tFlow.Class,
	}

	for _, r := range tFlow.StartRules {
		newFlow.StartRules = append(newFlow.StartRules, rules{Model: gorm.Model{ID: r}})
	}
	for _, r := range tFlow.AllRules {
		newFlow.AllRules = append(newFlow.AllRules, rules{Model: gorm.Model{ID: r}})
	}

	if err := p.client.Create(&newFlow).Error; err != nil {
		return fmt.Errorf("method *PostgresApiRepository.InsertTriggerFlow: could not insert trigger flows: %s", err)
	}

	return nil
}
