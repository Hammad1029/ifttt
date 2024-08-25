package main

import (
	"fmt"
	"ifttt/manager/application/config"
	"ifttt/manager/application/server"
)

func main() {
	fmt.Println("Starting application")
	if err := config.Init(); err != nil {
		panic(fmt.Errorf("could not init config %s", err))
	}
	if err := server.Init(); err != nil {
		panic(fmt.Errorf("could not init server %s", err))
	}
}
