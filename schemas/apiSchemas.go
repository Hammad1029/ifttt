package schemas

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gocql/gocql"
)

type rule struct {
	Id        gocql.UUID `json:"id"`
	Operator1 string     `json:"op1"`
	Operand   string     `json:"opnd"`
	Operator2 string     `json:"op2"`
	Then      []action   `json:"then"`
	Else      []action   `json:"else"`
}

type action struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type AddApiRequest struct {
	ApiGroup       string `json:"apiGroup"`
	ApiName        string `json:"apiName"`
	ApiDescription string `json:"apiDescription"`
	ApiPath        string `json:"apiPath"`
	Rules          []rule `json:"rules"`
}

func (r AddApiRequest) Validate() error {
	// operators := schemas.GetStringSlice("operators")
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.ApiGroup,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.ApiName,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.ApiDescription,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.ApiPath,
			validation.Required,
			validation.Length(3, 0),
			// validation.Match(regexp.MustCompile(utils.Regex.Endpoint)),
		),
		validation.Field(
			&r.Rules,
			validation.Required,
			validation.Length(1, 100),
			validation.Each(
			// here goes the validation for rules
			),
		),
	)
}

/*
validation.Field(
	(*rule).Operator1,
	validation.Required,
	validation.Length(1, 0),
),
validation.Field(
	(*rule).Operand,
	validation.Required,
	validation.Each(validation.In(utils.ConvertStringToInterfaceArray(operators)...)),
),
validation.Field(
	(*rule).Operator2,
	validation.Required,
	validation.Length(1, 0),
),
validation.Field(
	(*rule).Id,
	validation.Required,
	validation.Match(regexp.MustCompile(utils.Regex.UUID)),
),
validation.Field(
	(*rule).Then,
	validation.Required,
),
validation.Field(
	(*rule).Else,
	validation.Required,
),
*/
