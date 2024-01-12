package utils

import (
	"generic/config"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Roles = map[string]int{
	"unauthorized": 0,
	"authorized":   1,
	"admin":        2,
}

var CustomValidator *validator.Validate

func init() {
	CustomValidator = validator.New()
}

func Validator(controller func(*gin.Context), authorized int, schema interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authorized == Roles["authorized"] || authorized == Roles["admin"] {
			tokenString := c.GetHeader("Authorization")

			if tokenString == "" {
				ResponseHandler(c, Config{Response: Responses["Unauthorized"]})
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.GetConfigProp("jwt.secret")), nil
			})

			if err != nil || !token.Valid {
				ResponseHandler(c, Config{Response: Responses["Unauthorized"]})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				ResponseHandler(c, Config{Response: Responses["Unauthorized"]})
				return
			}

			if authorized == Roles["admin"] && !claims["admin"].(bool) {
				ResponseHandler(c, Config{Response: Responses["Unauthorized"]})
				return
			}

			c.Set("user", claims)
		}

		if schema != nil {
			schemaType := reflect.TypeOf(schema)
			schemaInstance := reflect.New(schemaType).Interface()
			if err := c.ShouldBindJSON(&schemaInstance); err != nil {
				ValidationErrors := make(map[string]interface{})
				ValidationErrors["Validation Errors"] = strings.Split(err.Error(), "\n")
				ResponseHandler(c, Config{Response: Responses["BadRequest"], Data: ValidationErrors})
				return
			}
			c.Set("Req", schemaInstance)
		}

		controller(c)
	}
}
