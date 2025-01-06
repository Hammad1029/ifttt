package postgres

import (
	"fmt"
	"ifttt/manager/domain/api"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"gorm.io/gorm"
)

type PostgresAPIRepository struct {
	*PostgresBaseRepository
}

func NewPostgresAPIRepository(base *PostgresBaseRepository) *PostgresAPIRepository {
	return &PostgresAPIRepository{PostgresBaseRepository: base}
}

func (p *PostgresAPIRepository) GetAllApis() (*[]api.Api, error) {
	var pgApis []apis
	if err := p.client.
		Preload("TriggerFlowRef").Preload("TriggerFlowRef.Rules").
		Find(&pgApis).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not query apis: %s", err)
	}

	var domainApis []api.Api
	for _, a := range pgApis {
		if dApi, err := a.toDomain(); err != nil {
			return nil,
				fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not convert to domain api: %s", err)
		} else {
			domainApis = append(domainApis, *dApi)
		}
	}

	return &domainApis, nil
}

func (p *PostgresAPIRepository) GetApiByNameOrPath(name string, path string) (*api.Api, error) {
	var pgApi apis
	if err := p.client.
		Preload("Triggers").Preload("Triggers.Rules").
		First(&pgApi, "name = ? or path = ?", name, path).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("method*PostgresApiRepository.GetApiDetails: could not query apis: %s", err)
	}

	var domainApi api.Api
	if dApi, err := pgApi.toDomain(); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not convert to domain api: %s", err)
	} else {
		domainApi = *dApi
	}

	return &domainApi, nil
}

func (p *PostgresAPIRepository) InsertApi(apiReq *api.CreateApiRequest, attachTriggers *[]triggerflow.TriggerFlow) error {
	var pgApi apis
	if err := pgApi.fromDomain(apiReq, attachTriggers); err != nil {
		return fmt.Errorf("could not convert from domain api: %s", err)
	}

	if err := p.client.Create(&pgApi).Error; err != nil {
		return fmt.Errorf("could not insert api: %s", err)
	}
	return nil
}
