package controllers

import infrastructure "ifttt/manager/infrastructure/init"

type AllController struct {
	ApiController *apiController
}

func NewAllController(store *infrastructure.DbStore) *AllController {
	return &AllController{
		ApiController: newApiController(store),
	}
}
