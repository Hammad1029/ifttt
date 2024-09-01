package middlewares

import "ifttt/manager/application/core"

type allMiddlewares struct {
	serverCore *core.ServerCore
}

func NewAllMiddlewares(serverCore *core.ServerCore) *allMiddlewares {
	return &allMiddlewares{
		serverCore: serverCore,
	}
}
