package config

import "rogchap.com/v8go"

var v8Isolate *v8go.Isolate
var v8Ctx *v8go.Context

func startV8() {
	v8Isolate = v8go.NewIsolate()
	defer v8Isolate.Dispose()
	v8Ctx = v8go.NewContext(v8Isolate)
	defer v8Ctx.Close()
}

func GetV8() *v8go.Context {
	return v8Ctx
}
