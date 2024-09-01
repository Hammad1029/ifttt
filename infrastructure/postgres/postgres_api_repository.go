package postgres

import (
	"fmt"
	"ifttt/manager/domain/api"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type postgresAPI struct {
	gorm.Model
	Group       string         `gorm:"type:varchar(50)" mapstructure:"group"`
	Name        string         `gorm:"type:varchar(50)" mapstructure:"name"`
	Method      string         `gorm:"type:varchar(10)" mapstructure:"method"`
	Type        string         `gorm:"type:varchar(10)" mapstructure:"type"`
	Path        string         `gorm:"type:varchar(50)" mapstructure:"path"`
	Description string         `gorm:"type:varchar(255)" mapstructure:"description"`
	Request     map[string]any `gorm:"type:jsonb;default:'{}';not null" mapstructure:"request"`
	Dumping     map[string]any `gorm:"type:jsonb;default:'{}';not null" mapstructure:"dumping"`
	StartRules  []string       `gorm:"type:varchar(50)[];default:'[]';not null" mapstructure:"rules"`
	Rules       map[string]any `gorm:"type:jsonb;default:'{}';not null" mapstructure:"startRules"`
}

func (p postgresAPI) TableName() string {
	return "apis"
}

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
