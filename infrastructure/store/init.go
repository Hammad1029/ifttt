package store

import (
	"fmt"
	"ifttt/manager/application/config"
	"ifttt/manager/domain/api"
	"ifttt/manager/domain/schema"
	"ifttt/manager/domain/token"
	"ifttt/manager/domain/user"
	"strings"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter"
)

type dbStorer interface {
	init(config map[string]any) error
}

type configStorer interface {
	dbStorer
	createConfigStore() *ConfigStore
	createCasbinAdapter() *gormadapter.Adapter
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
	Store          configStorer
	CasbinEnforcer *casbin.Enforcer
	APIRepo        api.Repository
	UserRepo       user.Repository
}

type DataStore struct {
	Store      dataStorer
	SchemaRepo schema.Repository
}

type CacheStore struct {
	Store     cacheStorer
	TokenRepo token.Repository
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
	// case scyllaDb:
	// 	storer = &scyllaStore{}
	case postgresDb:
		storer = &postgresStore{}
	default:
		return nil, fmt.Errorf("method configStoreFactory: db not found %s", dbName)
	}

	if err := storer.init(connectionSettings); err != nil {
		return nil, fmt.Errorf("method configStoreFactory: could not init config store: %s", err)
	}
	configStore := storer.createConfigStore()

	casbinAdapter := storer.createCasbinAdapter()
	casbinEnforcer, err := casbin.NewEnforcer("./application/config/casbin_model.conf", casbinAdapter)
	if err != nil {
		return nil, fmt.Errorf("method configStoreFactory: could not create casbin enforcer: %s", err)
	}
	if err := casbinEnforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("method configStoreFactory: could not load casbin policy: %s", err)
	}
	configStore.CasbinEnforcer = casbinEnforcer

	return configStore, nil
}

func dataStoreFactory(connectionSettings map[string]any) (*DataStore, error) {
	var storer dataStorer
	dbName, ok := connectionSettings["db"]
	if !ok {
		return nil, fmt.Errorf("method dataStoreFactory: db name not found in env")
	}

	switch strings.ToLower(fmt.Sprint(dbName)) {
	case scyllaDb:
		storer = &scyllaStore{}
	case postgresDb:
		storer = &postgresStore{}
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
	case redisCache:
		storer = &RedisStore{}
	default:
		return nil, fmt.Errorf("method cacheStoreFactory: db not found %s", dbName)
	}

	if err := storer.init(connectionSettings); err != nil {
		return nil, fmt.Errorf("method configStoreFactory: could not init cache store: %s", err)
	}

	return storer.createCacheStore(), nil
}
