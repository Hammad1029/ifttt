package common

import (
	"log"

	"github.com/gin-gonic/gin"
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

func HandleErrorResponse(c *gin.Context, e any, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Println(msg[0], e)
		}
		ResponseHandler(c, ResponseConfig{Response: Responses["ServerError"]})
		panic(e)
	}
}
