package scylla

import (
	"fmt"
	"generic/domain/api"

	"github.com/mitchellh/mapstructure"
	"github.com/scylladb/gocqlx/v3/table"
)

type scyllaApiSerialized struct {
	ApiGroup       string   `cql:"api_group"`
	ApiName        string   `cql:"api_name"`
	ApiDescription string   `cql:"api_description"`
	ApiPath        string   `cql:"api_path"`
	ApiRequest     string   `cql:"api_request"`
	Rules          string   `cql:"rules"`
	StartRules     []string `cql:"start_rules"`
}

var scyllaApisMetadata = table.Metadata{
	Name:    "apis",
	Columns: []string{"api_group", "api_name", "api_description", "api_path", "api_request", "rules", "start_rules"},
	PartKey: []string{"api_group"},
	SortKey: []string{"api_name", "api_description"},
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

func (s *ScyllaApiRepository) GetAllApis() (*[]api.ApiSerialized, error) {
	var scyllaApis []scyllaApiSerialized
	serializedApis := &([]api.ApiSerialized{})

	apisTable := s.getTable()
	stmt, names := apisTable.SelectAll()
	if err := s.session.Query(stmt, names).SelectRelease(&scyllaApis); err != nil {
		return nil, fmt.Errorf("method *ScyllaApiRepository.GetAllApis: could not get apis: %s", err)
	}

	if err := mapstructure.Decode(scyllaApis, &serializedApis); err != nil {
		return nil, fmt.Errorf("method *ScyllaApiRepository.GetAllApis: failed to decode apis: %s", err)
	}

	return serializedApis, nil
}

func (s *ScyllaApiRepository) GetApiByGroupAndName(group string, name string) (*api.ApiSerialized, bool, error) {
	var scyllaApis []scyllaApiSerialized
	serializdApi := &(api.ApiSerialized{})
	var found bool

	apisTable := s.getTable()
	query := apisTable.SelectQuery(*s.session).BindStruct(scyllaApiSerialized{
		ApiGroup: group,
		ApiName:  name,
	})
	if err := query.SelectRelease(&scyllaApis); err != nil {
		return nil, found, fmt.Errorf("method *ScyllaApiRepository.GetApiByGroupAndName: failed to get api: %s", err)
	}

	if len(scyllaApis) == 0 {
		return nil, found, nil
	}
	found = true

	if err := mapstructure.Decode(scyllaApis[0], &serializdApi); err != nil {
		return nil, found, fmt.Errorf("method *ScyllaApiRepository.GetApiByGroupAndName: failed to decode api: %s", err)
	}

	return serializdApi, found, nil
}

func (s *ScyllaApiRepository) InsertApi(newApi *api.ApiSerialized) error {
	var scyllaApiSerialized scyllaApiSerialized

	if err := mapstructure.Decode(newApi, &scyllaApiSerialized); err != nil {
		return fmt.Errorf("method *ScyllaApiRepository.InsertApi: failed to decode api: %s", err)
	}

	apisTable := s.getTable()
	query := apisTable.InsertQuery(*s.session).BindStruct(scyllaApiSerialized)
	if err := query.ExecRelease(); err != nil {
		return fmt.Errorf("method *ScyllaApiRepository.InsertApi: failed to insert api: %s", err)
	}

	return nil
}
