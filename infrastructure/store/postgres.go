package store

import (
	"fmt"
	postgresInfra "ifttt/manager/infrastructure/postgres"

	"github.com/mitchellh/mapstructure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const postgresDb = "postgres"

type postgresStore struct {
	store  *gorm.DB
	config postgresConfig
}

type postgresConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	Database string `json:"database" mapstructure:"database"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

func (p *postgresStore) init(config map[string]any) error {
	if err := mapstructure.Decode(config, &p.config); err != nil {
		return fmt.Errorf(
			"method: *PostgresStore.Init: could not decode scylla configuration from env: %s", err,
		)
	}
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Karachi",
		p.config.Host, p.config.Username, p.config.Password, p.config.Database, p.config.Port,
	)
	if db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{}); err != nil {
		return err
	} else {
		p.store = db
	}
	return nil
}

func (p *postgresStore) createConfigStore() *ConfigStore {
	postgresBase := postgresInfra.NewPostgresBaseRepository(p.store)
	return &ConfigStore{
		Store:    p,
		APIRepo:  postgresInfra.NewPostgresAPIRepository(postgresBase),
		UserRepo: postgresInfra.NewPostgresUserRepository(postgresBase),
	}
}

func (p *postgresStore) createDataStore() *DataStore {
	postgresBase := postgresInfra.NewPostgresBaseRepository(p.store)
	return &DataStore{
		Store:      p,
		SchemaRepo: postgresInfra.NewPostgresSchemaRepository(postgresBase),
	}
}
