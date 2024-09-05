package core

import (
	"fmt"
	"ifttt/manager/domain/token"
	"ifttt/manager/infrastructure/store"
)

type ServerCore struct {
	ConfigStore  *store.ConfigStore
	DataStore    *store.DataStore
	CacheStore   *store.CacheStore
	TokenService *token.TokenService
}

func NewServerCore() (*ServerCore, error) {
	var serverCore ServerCore

	if configStore, err := store.NewConfigStore(); err != nil {
		return nil, fmt.Errorf("method NewServerCore: could not create config store: %s", err)
	} else {
		serverCore.ConfigStore = configStore
	}
	if dataStore, err := store.NewDataStore(); err != nil {
		return nil, fmt.Errorf("method NewServerCore: could not create data store: %s", err)
	} else {
		serverCore.DataStore = dataStore
	}
	if cacheStore, err := store.NewCacheStore(); err != nil {
		return nil, fmt.Errorf("method NewServerCore: could not create cache store: %s", err)
	} else {
		serverCore.CacheStore = cacheStore
	}
	if tokenService, err := token.NewTokenService(); err != nil {
		return nil, fmt.Errorf("method NewServerCore: could not create token service: %s", err)
	} else {
		serverCore.TokenService = tokenService
	}

	return &serverCore, nil
}
