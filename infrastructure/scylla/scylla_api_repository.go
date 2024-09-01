package scylla

import (
	"encoding/json"
	"fmt"
	"ifttt/manager/domain/api"
	"reflect"

	"github.com/scylladb/gocqlx/v3/table"
)

type localAPI struct {
	Group       string   `cql:"group" mapstructure:"group"`
	Name        string   `cql:"name" mapstructure:"name"`
	Method      string   `cql:"method" mapstructure:"method"`
	Type        string   `cql:"type" mapstructure:"type"`
	Path        string   `cql:"path" mapstructure:"path"`
	Description string   `cql:"description" mapstructure:"description"`
	Request     string   `cql:"request" mapstructure:"request"`
	Dumping     string   `cql:"dumping" mapstructure:"dumping"`
	StartRules  []string `cql:"start_rules" mapstructure:"rules"`
	Rules       string   `cql:"rules" mapstructure:"startRules"`
}

var scyllaApisMetadata = table.Metadata{
	Name:    "apis",
	Columns: []string{"group", "name", "method", "type", "path", "description", "request", "dumping", "start_rules", "rules"},
	PartKey: []string{"group"},
	SortKey: []string{"name"},
}

var scyllaApisTable *table.Table

type ScyllaApiRepository struct {
	ScyllaBaseRepository
}

func NewScyllaApiRepository(base ScyllaBaseRepository) *ScyllaApiRepository {
	return &ScyllaApiRepository{ScyllaBaseRepository: base}
}

func (s *ScyllaApiRepository) getTable() *table.Table {
	if scyllaApisTable == nil {
		scyllaApisTable = table.New(scyllaApisMetadata)
	}
	return scyllaApisTable
}

func (s *ScyllaApiRepository) ToLocal(input *api.Api, output any) error {
	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return fmt.Errorf("method *ScyllaApiRepository: output not pointer")
	}
	serializedAPI := localAPI{
		Group:       (*input).Group,
		Name:        (*input).Name,
		Method:      (*input).Method,
		Type:        (*input).Type,
		Path:        (*input).Path,
		Description: (*input).Description,
		StartRules:  (*input).StartRules,
	}

	requestSerialized, err := json.Marshal((*input).Request)
	if err != nil {
		return fmt.Errorf("method *Api.TransformApiForSave: could not serialize request: %s", err)
	}
	serializedAPI.Request = string(requestSerialized)

	rulesSerialized, err := json.Marshal((*input).Rules)
	if err != nil {
		return fmt.Errorf("method *Api.TransformApiForSave: could not serialize rules: %s", err)
	}
	serializedAPI.Rules = string(rulesSerialized)

	dumpingSerialized, err := json.Marshal((*input).Dumping)
	if err != nil {
		return fmt.Errorf("method *Api.TransformApiForSave: could not serialize dumping: %s", err)
	}
	serializedAPI.Dumping = string(dumpingSerialized)

	output = &serializedAPI

	return nil
}

func (s *ScyllaApiRepository) ToGlobal(input any) (*api.Api, error) {
	localAPI, ok := input.(localAPI)
	if !ok {
		return nil, fmt.Errorf("method *ScyllaApiRepository.ToGlobal: could not convert api to local type")
	}

	globalAPI := api.Api{
		Group:       localAPI.Group,
		Description: localAPI.Description,
		Name:        localAPI.Name,
		Type:        localAPI.Type,
		Path:        localAPI.Path,
		Method:      localAPI.Method,
		StartRules:  localAPI.StartRules,
	}

	if err := json.Unmarshal([]byte(localAPI.Request), &globalAPI.Request); err != nil {
		return nil, fmt.Errorf("method *ScyllaApiRepository.ToGlobal: could not unmarshal request: %s", err)
	}

	if err := json.Unmarshal([]byte(localAPI.Dumping), &globalAPI.Dumping); err != nil {
		return nil, fmt.Errorf("method *ScyllaApiRepository.ToGlobal: could not unmarshal dumping: %s", err)
	}

	if err := json.Unmarshal([]byte(localAPI.Rules), &globalAPI.Rules); err != nil {
		return nil, fmt.Errorf("method *ScyllaApiRepository.ToGlobal: could not unmarshal request: %s", err)
	}

	return &globalAPI, nil
}

func (s *ScyllaApiRepository) GetAllApis() (*[]api.Api, error) {
	var scyllaApis []localAPI
	var globalApis []api.Api

	apisTable := s.getTable()
	stmt, names := apisTable.SelectAll()
	if err := s.session.Query(stmt, names).SelectRelease(&scyllaApis); err != nil {
		return nil, fmt.Errorf("method *ScyllaApiRepository.GetAllApis: could not get apis: %s", err)
	}

	for _, localApi := range scyllaApis {
		if global, err := s.ToGlobal(localApi); err != nil {
			return nil, fmt.Errorf("method *ScyllaApiRepository.GetAllApis: could not get global api: %s", err)
		} else {
			globalApis = append(globalApis, *global)
		}
	}

	return &globalApis, nil
}

func (s *ScyllaApiRepository) GetApisByGroupAndName(group string, name string) (*[]api.Api, error) {
	var globalApis []api.Api
	var scyllaApis []localAPI

	apisTable := s.getTable()
	query := apisTable.SelectQuery(*s.session).BindStruct(localAPI{
		Group: group,
		Name:  name,
	})
	if err := query.SelectRelease(&scyllaApis); err != nil {
		return nil, fmt.Errorf("method *ScyllaApiRepository.GetApisByGroupAndName: failed to get api: %s", err)
	}

	if len(scyllaApis) == 0 {
		return nil, nil
	}

	for _, localApi := range scyllaApis {
		if global, err := s.ToGlobal(localApi); err != nil {
			return nil, fmt.Errorf("method *ScyllaApiRepository.GetAllApis: could not get global api: %s", err)
		} else {
			globalApis = append(globalApis, *global)
		}
	}

	return &globalApis, nil
}

func (s *ScyllaApiRepository) InsertApi(newApi *api.Api) error {
	var serializedApi localAPI
	if err := s.ToLocal(newApi, &serializedApi); err != nil {
		return fmt.Errorf("method *ScyllaApiRepository.InsertApi: failed to convert api to local struct: %s", err)
	}

	apisTable := s.getTable()
	query := apisTable.InsertQuery(*s.session).BindStruct(serializedApi)
	if err := query.ExecRelease(); err != nil {
		return fmt.Errorf("method *ScyllaApiRepository.InsertApi: failed to insert api: %s", err)
	}

	return nil
}
