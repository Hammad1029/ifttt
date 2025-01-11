package store

import (
	"fmt"
	"ifttt/manager/application/config"
	"ifttt/manager/common"
	"ifttt/manager/domain/api"
	"ifttt/manager/domain/auth"
	"ifttt/manager/domain/cron"
	"ifttt/manager/domain/orm_schema"
	"ifttt/manager/domain/resolvable"
	responseprofiles "ifttt/manager/domain/response_profiles"
	"ifttt/manager/domain/rule"
	triggerflow "ifttt/manager/domain/trigger_flow"
	"ifttt/manager/domain/user"
	"strings"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type dbStorer interface {
	init(config map[string]any) error
}

type configStorer interface {
	dbStorer
	createConfigStore() *ConfigStore
	createCasbinAdapter() (*gormadapter.Adapter, error)
}

type dataStorer interface {
	dbStorer
	createDataStore() *DataStore
}

type cacheStorer interface {
	dbStorer
	createCacheStore() *CacheStore
}

type ConfigStore struct {
	Store               configStorer
	CasbinEnforcer      *casbin.Enforcer
	CronRepo            cron.Repository
	APIRepo             api.Repository
	RuleRepo            rule.Repository
	TriggerFlowRepo     triggerflow.Repository
	UserRepo            user.Repository
	OrmRepo             orm_schema.OrmRepository
	ResponseProfileRepo responseprofiles.Repository
}

type DataStore struct {
	Store                 dataStorer
	SchemaRepo            orm_schema.SchemaRepository
	OrmQueryGeneratorRepo resolvable.OrmQueryGenerator
}

type CacheStore struct {
	Store    cacheStorer
	AuthRepo auth.Repository
}

func NewConfigStore() (*ConfigStore, error) {
	connectionSettings := config.GetConfig().GetStringMap("configStore")
	if store, err := configStoreFactory(connectionSettings); err != nil {
		return nil, fmt.Errorf("method NewConfigStore: could not create store: %s", err)
	} else {
		return store, nil
	}
}

func NewDataStore() (*DataStore, error) {
	connectionSettings := config.GetConfig().GetStringMap("dataStore")
	if store, err := dataStoreFactory(connectionSettings); err != nil {
		return nil, fmt.Errorf("method NewDataStore: could not create store: %s", err)
	} else {
		return store, nil
	}
}

func NewCacheStore() (*CacheStore, error) {
	connectionSettings := config.GetConfig().GetStringMap("cacheStore")
	if store, err := cacheStoreFactory(connectionSettings); err != nil {
		return nil, fmt.Errorf("method InitCacheStore: could not create store: %s", err)
	} else {
		return store, nil
	}
}

func configStoreFactory(connectionSettings map[string]any) (*ConfigStore, error) {
	var storer configStorer
	dbName, ok := connectionSettings["db"]
	if !ok {
		return nil, fmt.Errorf("method configStoreFactory: db name not found in env")
	}

	switch strings.ToLower(fmt.Sprint(dbName)) {
	case common.DbNamePostgres:
		storer = &postgresStore{}
	default:
		return nil, fmt.Errorf("method configStoreFactory: db not found %s", dbName)
	}

	if err := storer.init(connectionSettings); err != nil {
		return nil, fmt.Errorf("method configStoreFactory: could not init config store: %s", err)
	}
	configStore := storer.createConfigStore()

	if err := casbinFactory(storer, configStore); err != nil {
		return nil, fmt.Errorf("method configStoreFactory: could not init casbin: %s", err)
	}

	return configStore, nil
}

func dataStoreFactory(connectionSettings map[string]any) (*DataStore, error) {
	var storer dataStorer
	dbName, ok := connectionSettings["db"]
	if !ok {
		return nil, fmt.Errorf("method dataStoreFactory: db name not found in env")
	}

	switch strings.ToLower(fmt.Sprint(dbName)) {
	case common.DbNamePostgres:
		storer = &postgresStore{}
	case common.DbNameMySql:
		storer = &mysqlStore{}
	default:
		return nil, fmt.Errorf("method dataStoreFactory: db not found %s", dbName)
	}

	if err := storer.init(connectionSettings); err != nil {
		return nil, fmt.Errorf("method dataStoreFactory: could not init data store: %s", err)
	}
	return storer.createDataStore(), nil
}

func cacheStoreFactory(connectionSettings map[string]any) (*CacheStore, error) {
	var storer cacheStorer
	dbName, ok := connectionSettings["db"]
	if !ok {
		return nil, fmt.Errorf("method cacheStoreFactory: db name not found in env")
	}

	switch strings.ToLower(fmt.Sprint(dbName)) {
	case common.DbNameRedis:
		storer = &RedisStore{}
	default:
		return nil, fmt.Errorf("method cacheStoreFactory: db not found %s", dbName)
	}

	if err := storer.init(connectionSettings); err != nil {
		return nil, fmt.Errorf("method configStoreFactory: could not init cache store: %s", err)
	}

	return storer.createCacheStore(), nil
}

func casbinFactory(storer configStorer, store *ConfigStore) error {
	casbinAdapter, err := storer.createCasbinAdapter()
	if err != nil {
		return fmt.Errorf("method casbinFactory: could not create casbin adapter: %s", err)
	}

	casbinEnforcer, err := casbin.NewEnforcer(config.GetConfig().GetString("casbin.modelLocation"), casbinAdapter)
	if err != nil {
		return fmt.Errorf("method casbinFactory: could not create casbin enforcer: %s", err)
	}
	if err := casbinEnforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("method casbinFactory: could not load casbin policy: %s", err)
	}
	store.CasbinEnforcer = casbinEnforcer

	// if auth, err := gcasbin.NewCasbinMiddleware(
	// 	"./application/config/casbin_model.conf", casbinAdapter, user.GetEmailFromContext); err != nil {
	// 	return fmt.Errorf("method casbinFactory: could not create casbin enforcer: %s", err)
	// } else {
	// 	store.CasbinMiddleware = auth
	// }

	return nil
}
