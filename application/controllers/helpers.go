package controllers

import (
	"fmt"
	"ifttt/manager/common"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type validatorInterface interface {
	Validate() error
}

func validateRequest(c *gin.Context, output any) bool {
	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		common.HandleErrorResponse(c, fmt.Errorf("method validateRequest: output struct is not a pointer"))
		return false
	}

	if err := c.ShouldBindJSON(output); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return false
	}

	// if validator, ok := output.(validatorInterface); ok {
	// if err := validator.Validate(); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, err)
	// 	return false
	// }
	// }

	return true
}
