package store

import (
	"fmt"
	postgresInfra "ifttt/manager/infrastructure/postgres"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-viper/mapstructure/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresStore struct {
	store  *gorm.DB
	config postgresConfig
}

type postgresConfig struct {
	Host             string `json:"host" mapstructure:"host"`
	Port             string `json:"port" mapstructure:"port"`
	Database         string `json:"database" mapstructure:"database"`
	Username         string `json:"username" mapstructure:"username"`
	Password         string `json:"password" mapstructure:"password"`
	ConnectionString string `json:"connectionString" mapstructure:"connectionString"`
}

func (p *postgresStore) init(config map[string]any) error {
	if err := mapstructure.Decode(config, &p.config); err != nil {
		return fmt.Errorf(
			"method: *PostgresStore.Init: could not decode scylla configuration from env: %s", err,
		)
	}
	p.config.ConnectionString = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		p.config.Host, p.config.Username, p.config.Password, p.config.Database, p.config.Port,
	)
	if db, err := gorm.Open(postgres.Open(p.config.ConnectionString), &gorm.Config{}); err != nil {
		return err
	} else {
		p.store = db
	}
	return nil
}

func (p *postgresStore) createConfigStore() *ConfigStore {
	postgresBase := postgresInfra.NewPostgresBaseRepository(p.store, true)
	return &ConfigStore{
		Store:               p,
		CronRepo:            postgresInfra.NewPostgresCronRepository(postgresBase),
		APIRepo:             postgresInfra.NewPostgresAPIRepository(postgresBase),
		RuleRepo:            postgresInfra.NewPostgresRulesRepository(postgresBase),
		TriggerFlowRepo:     postgresInfra.NewPostgresTriggerFlowsRepository(postgresBase),
		UserRepo:            postgresInfra.NewPostgresUserRepository(postgresBase),
		OrmRepo:             postgresInfra.NewPostgresOrmRepository(postgresBase),
		ResponseProfileRepo: postgresInfra.NewPostgresResponseProfileRepository(postgresBase),
		InternalTagRepo:     postgresInfra.NewPostgresInternalTagRepository(postgresBase),
	}
}

func (p *postgresStore) createCasbinAdapter() (*gormadapter.Adapter, error) {
	return gormadapter.NewAdapterByDB(p.store)
}

func (p *postgresStore) createDataStore() *DataStore {
	postgresBase := postgresInfra.NewPostgresBaseRepository(p.store, false)
	return &DataStore{
		Store:                 p,
		SchemaRepo:            postgresInfra.NewPostgresSchemaRepository(postgresBase),
		OrmQueryGeneratorRepo: postgresInfra.NewPostgresOrmQueryGeneratorRepository(postgresBase),
	}
}
