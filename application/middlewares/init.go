package middlewares

import "ifttt/manager/application/core"

type AllMiddlewares struct {
	serverCore *core.ServerCore
}

func NewAllMiddlewares(serverCore *core.ServerCore) *AllMiddlewares {
	return &AllMiddlewares{
		serverCore: serverCore,
	}
}
