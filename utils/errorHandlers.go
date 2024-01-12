package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, e any, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Println(msg[0], e)
		}
		ResponseHandler(c, Config{Response: Responses["ServerError"]})
		panic(e)
	}
}
