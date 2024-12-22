package store

import (
	"fmt"
	mysqlInfra "ifttt/manager/infrastructure/mysql"

	"github.com/mitchellh/mapstructure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlStore struct {
	store  *gorm.DB
	config mysqlConfig
}

type mysqlConfig struct {
	Host             string `json:"host" mapstructure:"host"`
	Port             string `json:"port" mapstructure:"port"`
	Database         string `json:"database" mapstructure:"database"`
	Username         string `json:"username" mapstructure:"username"`
	Password         string `json:"password" mapstructure:"password"`
	ConnectionString string `json:"connectionString" mapstructure:"connectionString"`
}

func (m *mysqlStore) init(config map[string]any) error {
	if err := mapstructure.Decode(config, &m.config); err != nil {
		return fmt.Errorf("method: *mysqlStore.Init: could not decode configuration from env: %s", err)
	}
	m.config.ConnectionString = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		m.config.Username, m.config.Password, m.config.Host, m.config.Port, m.config.Database,
	)
	if db, err := gorm.Open(mysql.Open(m.config.ConnectionString), &gorm.Config{}); err != nil {
		return err
	} else {
		m.store = db
	}
	return nil
}

func (m *mysqlStore) createDataStore() *DataStore {
	mysqlBase := mysqlInfra.NewMySqlBaseRepository(m.store)
	return &DataStore{
		Store:      m,
		SchemaRepo: mysqlInfra.NewMySqlSchemaRepository(mysqlBase),
	}
}
