package postgres

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type users struct {
	gorm.Model
	Email    string `gorm:"type:varchar(50);unique" mapstructure:"email"`
	Password string `gorm:"type:varchar(255)" mapstructure:"password"`
}

type apis struct {
	gorm.Model
	Name         string          `gorm:"type:varchar(50);not null;unique" mapstructure:"name"`
	Path         string          `gorm:"type:varchar(50);not null;unique" mapstructure:"path"`
	Method       string          `gorm:"type:varchar(10);not null" mapstructure:"method"`
	Description  string          `gorm:"type:text;default:''" mapstructure:"description"`
	Request      pgtype.JSONB    `gorm:"type:jsonb;default:'{}';not null" mapstructure:"request"`
	PreConfig    pgtype.JSONB    `gorm:"type:jsonb;default:'{}';not null" mapstructure:"preConfig"`
	Triggerflows []trigger_flows `gorm:"many2many:api_trigger_flows;joinForeignKey:ApiId;joinReferences:FlowId;" mapstructure:"triggerFlows"`
}

type classes struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null;unique" mapstructure:"name"`
}

type trigger_flows struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(50);not null;unique" mapstructure:"name"`
	Description string  `gorm:"type:text;default:''" mapstructure:"description"`
	ClassId     uint    `gorm:"type:int;not null" mapstructure:"classId"`
	Class       classes `mapstructure:"class"`
	StartRules  []rules `gorm:"many2many:trigger_start_rules;joinForeignKey:FlowId;joinReferences:RuleId;" mapstructure:"startRules"`
	AllRules    []rules `gorm:"many2many:trigger_all_rules;joinForeignKey:FlowId;joinReferences:RuleId;" mapstructure:"allRules"`
}

type rules struct {
	gorm.Model
	Name        string       `json:"name" gorm:"type:varchar(50);not null;unique" mapstructure:"name"`
	Description string       `json:"description" gorm:"type:text;default:''" mapstructure:"description"`
	Conditions  pgtype.JSONB `json:"conditions" gorm:"type:jsonb;default:'{\"group\":true,\"conditionType\":\"and\",\"conditions\":[]}';not null" mapstructure:"conditions"`
	Then        pgtype.JSONB `json:"then" gorm:"type:jsonb;default:'[]';not null" mapstructure:"then"`
	Else        pgtype.JSONB `json:"else" gorm:"type:jsonb;default:'[]';not null" mapstructure:"else"`
}
