package postgres

import (
	"encoding/json"
	"fmt"
	"ifttt/manager/domain/api"

	"github.com/jackc/pgtype"
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
	if err := p.client.Find(&pgApis).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not query apis: %s", err)
	}

	var domainApis []api.Api
	for _, a := range pgApis {
		if dApi, err := p.ToDomain(&a); err != nil {
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
	if err := p.client.First(&pgApi, "name = ? or path = ?", name, path).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresAPIRepository.GetApisByNameOrPath: could not query apis: %s", err)
	}

	var domainApi api.Api
	if dApi, err := p.ToDomain(&pgApi); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not convert to domain api: %s", err)
	} else {
		domainApi = *dApi
	}

	return &domainApi, nil
}

func (p *PostgresAPIRepository) GetApiDetailsByNameAndPath(name string, path string) (*api.Api, error) {
	var pgApi apis
	if err := p.client.
		Preload("Triggerflows").Preload("Triggerflows.StartRules").Preload("Triggerflows.AllRules").
		First(&pgApi, "name = ? and path = ?", name, path).Error; err != nil {
		return nil, fmt.Errorf("method*PostgresApiRepository.GetApiDetails: could not query apis: %s", err)
	}

	var domainApi api.Api
	if dApi, err := p.ToDomain(&pgApi); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not convert to domain api: %s", err)
	} else {
		domainApi = *dApi
	}

	return &domainApi, nil
}

func (p *PostgresAPIRepository) InsertApi(apiReq *api.CreateApiRequest) error {
	newApi := api.Api{
		Name:        apiReq.Name,
		Path:        apiReq.Path,
		Method:      apiReq.Method,
		Description: apiReq.Name,
	}

	pgApi, err := p.FromDomain(&newApi)
	if err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not convert from domain api: %s", err)
	}
	pgApiCasted, ok := pgApi.(*apis)
	if !ok {
		return fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not cast pg api")
	}

	for _, a := range apiReq.TriggerFlows {
		pgApiCasted.Triggerflows = append(pgApiCasted.Triggerflows, trigger_flows{Model: gorm.Model{ID: a}})
	}

	if err := p.client.Create(pgApiCasted).Error; err != nil {
		return fmt.Errorf("method *PostgresApiRepository.InsertApi: could not insert api: %s", err)
	}
	return nil
}

func (p *PostgresAPIRepository) FromDomain(domainApi *api.Api) (any, error) {
	pgApi := apis{
		Name:        domainApi.Name,
		Path:        domainApi.Path,
		Method:      domainApi.Method,
		Description: domainApi.Description,
	}

	if reqMarshalled, err := json.Marshal(domainApi.Request); err != nil {
		return nil, fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal request")
	} else {
		pgApi.Request = pgtype.JSONB{Bytes: reqMarshalled}
	}

	if preConfigMarshalled, err := json.Marshal(domainApi.PreConfig); err != nil {
		return nil, fmt.Errorf("method *PostgresAPIRepository.FromDomain: could not marshal pre config")
	} else {
		pgApi.Request = pgtype.JSONB{Bytes: preConfigMarshalled}
	}

	return &pgApi, nil
}

func (p *PostgresAPIRepository) ToDomain(repoApi any) (*api.Api, error) {
	pgApi, ok := repoApi.(*apis)
	if !ok {
		return nil, fmt.Errorf("method *PostgresAPIRepository.ToDomain: could not cast pgApi")
	}

	domainApi := api.Api{
		Name:        pgApi.Name,
		Path:        pgApi.Path,
		Method:      pgApi.Method,
		Description: pgApi.Description,
	}

	if err := json.Unmarshal(pgApi.Request.Bytes, &domainApi.Request); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.ToDomain: could not cast pgApi: %s", err)
	}

	if err := json.Unmarshal(pgApi.PreConfig.Bytes, &domainApi.PreConfig); err != nil {
		return nil,
			fmt.Errorf("method *PostgresAPIRepository.ToDomain: could not cast pgApi: %s", err)
	}

	return &domainApi, nil
}
