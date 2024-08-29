package store

import (
	"fmt"
	scyllaInfra "ifttt/manager/infrastructure/scylla"
	"time"

	"github.com/gocql/gocql"
	"github.com/mitchellh/mapstructure"
	"github.com/scylladb/gocqlx/v3"
)

const scyllaDb = "scylla"

type scyllaStore struct {
	session *gocqlx.Session
	cluster *gocql.ClusterConfig
	config  scyllaConfig
}

type scyllaConfig struct {
	Keyspace   string   `json:"keyspace" mapstructure:"keyspace"`
	Nodes      []string `json:"nodes" mapstructure:"nodes"`
	Timeout    int      `json:"timeout" mapstructure:"timeout"`
	RetryMin   int      `json:"retryMin" mapstructure:"retryMin"`
	RetryMax   int      `json:"retryMax" mapstructure:"retryMax"`
	RetryCount int      `json:"retryCount" mapstructure:"retryCount"`
}

func (s *scyllaStore) init(config map[string]any) error {
	if err := mapstructure.Decode(config, &s.config); err != nil {
		return fmt.Errorf("method: *ScyllaStore.Init: could not decode scylla configuration from env: %s", err)
	}

	s.cluster = gocql.NewCluster(s.config.Nodes...)
	s.cluster.Keyspace = s.config.Keyspace
	s.cluster.Timeout = time.Duration(s.config.Timeout) * time.Second
	s.cluster.Consistency = gocql.Quorum
	s.cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	if session, err := gocqlx.WrapSession(gocql.NewSession(*s.cluster)); err != nil {
		return fmt.Errorf("method: *ScyllaStore.Init: error in creating new scylla session: %s", err)
	} else {
		s.session = &session
	}

	return nil
}

func (s *scyllaStore) createDataStore() *DataStore {
	// scyllaBase := scyllaInfra.NewScyllaBaseRepository(s.session, s.cluster)
	return &DataStore{
		Store: s,
	}
}

func (s *scyllaStore) createConfigStore() *ConfigStore {
	scyllaBase := scyllaInfra.NewScyllaBaseRepository(s.session, s.cluster)
	return &ConfigStore{
		Store:   s,
		APIRepo: scyllaInfra.NewScyllaApiRepository(*scyllaBase),
	}
}
