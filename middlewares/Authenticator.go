package middlewares

import (
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

func Authenticator() {
	// return func(c *gin.Context) {
	// 	if authorized == Roles["authorized"] || authorized == Roles["admin"] {
	// 		tokenString := c.GetHeader("Authorization")

	// 		if tokenString == "" {
	// 			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["Unauthorized"]})
	// 			return
	// 		}

	// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 			return []byte(config.GetConfigProp("jwt.secret")), nil
	// 		})

	// 		if err != nil || !token.Valid {
	// 			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["Unauthorized"]})
	// 			return
	// 		}

	// 		claims, ok := token.Claims.(jwt.MapClaims)
	// 		if !ok {
	// 			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["Unauthorized"]})
	// 			return
	// 		}

	// 		if authorized == Roles["admin"] && !claims["admin"].(bool) {
	// 			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["Unauthorized"]})
	// 			return
	// 		}

	// 		c.Set("user", claims)
	// 	}

	// 	if schema != nil {
	// 		schemaType := reflect.TypeOf(schema)
	// 		schemaInstance := reflect.New(schemaType).Interface()
	// 		if err := c.ShouldBindJSON(&schemaInstance); err != nil {
	// 			ValidationErrors := make(map[string]interface{})
	// 			ValidationErrors["Validation Errors"] = strings.Split(err.Error(), "\n")
	// 			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["BadRequest"], Data: ValidationErrors})
	// 			return
	// 		}
	// 		c.Set("Req", schemaInstance)
	// 	}

	// 	c.Next()
	// }
}
