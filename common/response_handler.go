package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	"Success":                  {"00", "Success"},
	"ClientNotFound":           {"05", "Client Not Found"},
	"ApiNotFound":              {"10", "Api Not Found"},
	"ApiAlreadyExists":         {"15", "API Already Exists"},
	"WrongTableFormat":         {"20", "Wrong Table Format"},
	"TableNotFound":            {"25", "Table Not Found"},
	"TableAlreadyExists":       {"25", "Table Already Exists"},
	"IndexNotPossible":         {"25", "Index Not Possible"},
	"IndexNotFound":            {"30", "Index Not Found"},
	"RuleAlreadyExists":        {"30", "Rule Already Exists"},
	"APIAlreadyExists":         {"30", "API Already Exists"},
	"FlowAlreadyExists":        {"30", "Trigger Flow Already Exists"},
	"TriggerFlowRulesNotFound": {"30", "All rules not found"},
	"TriggerFlowNotFound":      {"30", "Trigger flow not found"},
	"InvalidBranchFlow":        {"30", "Invalid Branch Flow"},
	"InvalidTriggerConditions": {"30", "Invalid Trigger Conditions"},
	"InvalidReturnValue":       {"30", "Invalid Return Value For Rule Case"},
	"CronAlreadyExists":        {"30", "Cron already exists"},

	"WrongCredentials":   {"35", "Wrong login credentials"},
	"UserNotFound":       {"40", "User not found"},
	"UserAlreadyExists":  {"40", "User already exists"},
	"PermissionNotFound": {"45", "Permission not found"},
	"RoleAlreadyExists":  {"45", "Role not found"},
	"BadRequest":         {"400", "Bad request"},
	"Unauthorized":       {"401", "Unauthorized"},
	"NotFound":           {"404", "Not Found"},
	"ServerError":        {"500", "Internal Server Error"},
}

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
	c.Abort()
}

func HandleErrorResponse(c *gin.Context, e any, msg ...string) {
	if e != nil {
		fmt.Print(e)
		ResponseHandler(c, ResponseConfig{Response: Responses["ServerError"]})
	}
}
