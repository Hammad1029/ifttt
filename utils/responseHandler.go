package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseHandler(c *gin.Context, config ResponseConfig) {

	if config.Response.Code == "" || config.Response.Description == "" {
		config.Response = Responses["Success"]
	}

	if config.Error != nil {
		HandleError(config.Error)
		config.Response = Responses["ServerError"]
	}

	if config.Data == nil {
		config.Data = make(map[string]any)
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

type ResponseConfig struct {
	Response Response
	Data     any
	Error    error
}

var Responses = map[string]Response{
	"Success":          {"00", "Success"},
	"ClientNotFound":   {"05", "Client Not Found"},
	"ApiNotFound":      {"10", "Api Not Found"},
	"ApiAlreadyExists": {"15", "API Already Exists"},
	"WrongTableFormat": {"20", "Wrong Table Format"},
	"TableNotFound":    {"25", "Table Not Found"},
	"IndexNotPossible": {"25", "Index Not Possible"},
	"IndexNotFound":    {"30", "Index Not Found"},

	"WrongCredentials": {"35", "Wrong login credentials"},
	"UserNotFound":     {"40", "User not found"},
	"BadRequest":       {"400", "Bad request"},
	"Unauthorized":     {"401", "Unauthorized"},
	"NotFound":         {"404", "Not Found"},
	"ServerError":      {"500", "Internal Server Error"},
}
