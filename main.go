package main

import (
	"fmt"
	"generic/application/config"
	"generic/application/server"
)

func main() {
	fmt.Println("Starting application")
	config.Init()
	server.Init()
}
