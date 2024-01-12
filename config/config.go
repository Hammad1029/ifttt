package config

import (
	"log"
)

func Init() {
	readEnv()
	startMongo()
	// startDruid()
	startV8()
	startKafka()
}

func handleError(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Println(msg[0], e)
		}
		panic(e)
	}
}
