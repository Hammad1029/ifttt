package controllers

import (
	"fmt"
	"ifttt/manager/common"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func validateAndBind(c *gin.Context, output common.Validatable) bool {
	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		common.HandleErrorResponse(c, fmt.Errorf("method validateAndBind: output struct is not a pointer"))
		return false
	}

	if err := c.ShouldBindJSON(output); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return false
	}

	if err := output.Validate(); err != nil {
		if internalErr, ok := err.(validation.InternalError); ok {
			common.HandleErrorResponse(c,
				fmt.Errorf("method validateAndBind: could not validate: %s", internalErr))
			return false
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return false
	}

	return true
}
