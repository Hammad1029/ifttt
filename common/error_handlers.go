package common

import (
	"log"
)

func HandleError(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Println(msg[0])
			log.Println(e)
		}
		panic(e)
	}
}
