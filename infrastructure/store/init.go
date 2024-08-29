package store

import (
	"fmt"
	"ifttt/manager/application/config"
	"ifttt/manager/domain/api"
	"ifttt/manager/domain/schema"
	"strings"
)

type dbStorer interface {
	init(config map[string]any) error
}

type configStorer interface {
	dbStorer
	createConfigStore() *ConfigStore
}

type dataStorer interface {
	dbStorer
	createDataStore() *DataStore
}

type ConfigStore struct {
	Store   configStorer
	APIRepo api.Repository
}

type DataStore struct {
	Store      dataStorer
	SchemaRepo schema.Repository
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

func configStoreFactory(connectionSettings map[string]any) (*ConfigStore, error) {
	var storer configStorer
	dbName, ok := connectionSettings["db"]
	if !ok {
		return nil, fmt.Errorf("method configStoreFactory: db name not found in env")
	}

	switch strings.ToLower(fmt.Sprint(dbName)) {
	case scyllaDb:
		storer = &scyllaStore{}
	default:
		return nil, fmt.Errorf("method configStoreFactory: db not found %s", dbName)
	}

	if err := storer.init(connectionSettings); err != nil {
		return nil, fmt.Errorf("method configStoreFactory: could not init config store: %s", err)
	}

	return storer.createConfigStore(), nil
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
