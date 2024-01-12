package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseHandler(c *gin.Context, params ...Config) {
	var config Config

	if len(params) > 0 {
		config = params[0]
	} else {
		config = Config{}
	}

	if config.Response.Code == "" || config.Response.Description == "" {
		config.Response = Responses["Success"]
	}

	if config.Error != nil {
		fmt.Println(config.Error)
		config.Response = Responses["ServerError"]
	}

	if config.Data == nil {
		config.Data = make(map[string]interface{})
	} else {
		config.Data = gin.H(config.Data)
	}

	c.JSON(http.StatusOK, gin.H{
		"responseCode":        config.Response.Code,
		"responseDescription": config.Response.Description,
		"data":                config.Data,
	})
}

type Response struct {
	Code        string
	Description string
}

type Config struct {
	Response Response
	Data     map[string]interface{}
	Error    error
}

var Responses = map[string]Response{
	"Success":          {"00", "Success"},
	"ClientNotFound":   {"05", "Client Not Found"},
	"ApiNotFound":      {"10", "Api Not Found"},
	"ApiAlreadyExists": {"15", "APIAlreadyExists"},

	"BadRequest":   {"400", "Bad request"},
	"Unauthorized": {"401", "Unauthorized"},
	"NotFound":     {"404", "Not Found"},
	"ServerError":  {"500", "Internal Server Error"},
}
