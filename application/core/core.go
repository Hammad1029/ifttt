package core

import (
	"fmt"
	"ifttt/manager/infrastructure/store"
)

type ServerCore struct {
	ConfigStore *store.ConfigStore
	DataStore   *store.DataStore
}

func NewServerCore() (*ServerCore, error) {
	var serverCore ServerCore

	if configStore, err := store.NewConfigStore(); err != nil {
		return nil, fmt.Errorf("method newCore: could not create config store: %s", err)
	} else {
		serverCore.ConfigStore = configStore
	}
	if dataStore, err := store.NewDataStore(); err != nil {
		return nil, fmt.Errorf("method newCore: could not create data store: %s", err)
	} else {
		serverCore.DataStore = dataStore
	}

	return &serverCore, nil
}
