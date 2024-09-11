package postgres

import (
	"fmt"
	"ifttt/manager/domain/api"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

type PostgresAPIRepository struct {
	*PostgresBaseRepository
}

func NewPostgresAPIRepository(base *PostgresBaseRepository) *PostgresAPIRepository {
	return &PostgresAPIRepository{PostgresBaseRepository: base}
}

func (p *PostgresAPIRepository) ToLocal(input *api.Api, output any) error {
	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return fmt.Errorf("method *ScyllaApiRepository: output not pointer")
	}

	var postgres postgresAPI
	if err := mapstructure.Decode((*input), &postgres); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.toLocal: could not decode to postgres: %s", err)
	}
	output = &postgres

	return nil
}

func (p *PostgresAPIRepository) ToGlobal(input any) (*api.Api, error) {
	postgresAPI, ok := input.(postgresAPI)
	if !ok {
		return nil, fmt.Errorf("method *PostgresAPIRepository.toGlobal: could not convert api to postgres type")
	}

	var globalApi api.Api
	if err := mapstructure.Decode(postgresAPI, &globalApi); err != nil {
		return nil, fmt.Errorf("method *PostgresAPIRepository.toGlobal: could not decode api: %s", err)
	}

	return &globalApi, nil
}

func (p *PostgresAPIRepository) GetAllApis() (*[]api.Api, error) {
	var globalApis []api.Api
	var postgresApis []postgresAPI
	if err := p.client.Find(&postgresApis).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresAPIRepository.GetAllApis: error in querying apis: %s", err)
	}

	for _, postgresApi := range postgresApis {
		if global, err := p.ToGlobal(postgresApi); err != nil {
			return nil, fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not get global api: %s", err)
		} else {
			globalApis = append(globalApis, *global)
		}
	}

	return &globalApis, nil
}

func (p *PostgresAPIRepository) GetApisByGroupAndName(group string, name string) (*[]api.Api, error) {
	var globalApis []api.Api
	var postgresApis []postgresAPI
	if err := p.client.Where(&postgresAPI{Group: group, Name: name}).Find(&postgresApis).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresAPIRepository.GetAllApis: error in querying apis: %s", err)
	}

	for _, postgresApi := range postgresApis {
		if global, err := p.ToGlobal(postgresApi); err != nil {
			return nil, fmt.Errorf("method *PostgresAPIRepository.GetAllApis: could not get global api: %s", err)
		} else {
			globalApis = append(globalApis, *global)
		}
	}

	return &globalApis, nil
}
func (p *PostgresAPIRepository) InsertApi(newApi *api.Api) error {
	var postgresAPI postgresAPI
	if err := p.ToLocal(newApi, &postgresAPI); err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.InsertApi: could not get postgres api: %s", err)
	}
	if err := p.client.Create(&postgresAPI).Error; err != nil {
		return fmt.Errorf("method *PostgresAPIRepository.InsertApi: could not insert api: %s", err)
	}

	return nil
}
