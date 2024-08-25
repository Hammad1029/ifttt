package infrastructure

import (
	"fmt"
	"ifttt/manager/application/config"
	"ifttt/manager/domain/api"
	"strings"
)

type dbStorer interface {
	init(config map[string]any) error
	createStore() *DbStore
}

type DbStore struct {
	Store         dbStorer
	ApiRepository api.Repository
}

func NewDbStore() (*DbStore, error) {
	connectionSettings := config.GetConfig().GetStringMap("store")
	if store, err := storeFactory(connectionSettings); err != nil {
		return nil, fmt.Errorf("method InitConfigStore: could not create store: %s", err)
	} else {
		return store, nil
	}
}

func storeFactory(connectionSettings map[string]any) (*DbStore, error) {
	var storer dbStorer
	dbName, ok := connectionSettings["db"]
	if !ok {
		return nil, fmt.Errorf("method storeFactory: db name not found in env")
	}

	switch strings.ToLower(fmt.Sprint(dbName)) {
	case scyllaDbName:
		storer = &scyllaStore{}
	}

	if err := storer.init(connectionSettings); err != nil {
		return nil, fmt.Errorf("method dataStoreFactory: could not init data store: %s", err)
	}
	return storer.createStore(), nil
}
