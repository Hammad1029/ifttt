package main

import (
	"fmt"

	"generic/config"
	"generic/server"
)

func main() {
	fmt.Println("Starting application")
	config.Init()
	server.Init()
}
