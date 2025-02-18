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

type crons struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50);not null;unique" mapstructure:"name"`
	Description string `gorm:"type:text;default:''" mapstructure:"description"`
	CronExpr    string `gorm:"type:varchar(30);default:''" mapstructure:"cronExpr"`
	ApiID       uint   `gorm:"not null" mapstructure:"apiId" json:"apiId"`
	API         apis   `gorm:"foreignKey:ApiID" mapstructure:"api" json:"api"`
}

type apis struct {
	gorm.Model
	Name         string          `gorm:"type:varchar(50);not null;unique" mapstructure:"name"`
	Path         string          `gorm:"type:varchar(50);not null;unique" mapstructure:"path"`
	Method       string          `gorm:"type:varchar(10);not null" mapstructure:"method"`
	Description  string          `gorm:"type:text;default:''" mapstructure:"description"`
	PreConfig    pgtype.JSONB    `gorm:"type:jsonb;default:'[]';not null" mapstructure:"preConfig"`
	Request      pgtype.JSONB    `gorm:"type:jsonb;default:'{}';not null" mapstructure:"request"`
	Response     pgtype.JSONB    `gorm:"type:jsonb;default:'{}';not null" mapstructure:"response"`
	Triggers     []trigger_flows `gorm:"many2many:api_trigger_flows_main;joinForeignKey:ApiId;joinReferences:FlowId;" mapstructure:"triggerFlows"`
	TriggerFlows pgtype.JSONB    `gorm:"type:jsonb;default:'{}';not null" mapstructure:"triggerConditions"`
}

type trigger_flows struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(50);not null;unique" mapstructure:"name"`
	Description string       `gorm:"type:text;default:''" mapstructure:"description"`
	StartState  uint         `gorm:"type:int;not null" mapstructure:"startState"`
	Rules       []rules      `gorm:"many2many:trigger_rules;joinForeignKey:FlowId;joinReferences:RuleId;" mapstructure:"rules"`
	BranchFlows pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null" mapstructure:"branchFlows"`
}

type rules struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(50);not null;unique" mapstructure:"name"`
	Description string       `gorm:"type:text;default:''" mapstructure:"description"`
	Pre         pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null" mapstructure:"pre"`
	Switch      pgtype.JSONB `gorm:"type:jsonb;default:'{\"cases\":[],\"default\":{\"do\":[],\"return\":{\"resolveType\":\"const\",\"resolveData\":\"\"}}}';not null" mapstructure:"switch"`
	Finally     pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null" mapstructure:"finally"`
}

type orm_model struct {
	gorm.Model
	Name                   string            `gorm:"type:varchar(255);not null" mapstructure:"name" json:"name"`
	Table                  string            `gorm:"type:varchar(255);not null" mapstructure:"table" json:"table"`
	PrimaryKey             string            `gorm:"type:varchar(255);not null" mapstructure:"primaryKey" json:"primaryKey"`
	Projections            []orm_projection  `gorm:"foreignKey:ModelID" mapstructure:"projections" json:"projections"`
	OwningAssociations     []orm_association `gorm:"foreignKey:OwningModelID" mapstructure:"owningAssociations" json:"owningAssociations"`
	ReferencedAssociations []orm_association `gorm:"foreignKey:ReferencesModelID" mapstructure:"referencedAssociations" json:"referencedAssociations"`
}

type orm_projection struct {
	gorm.Model
	ModelID    uint   `gorm:"not null"`
	Column     string `gorm:"type:varchar(255);not null" mapstructure:"column" json:"column"`
	As         string `gorm:"type:varchar(255);not null" mapstructure:"as" json:"as"`
	SchemaType string `gorm:"type:varchar(255);not null" mapstructure:"schemaType" json:"schemaType"`
	ModelType  string `gorm:"type:varchar(255);not null" mapstructure:"modelType" json:"modelType"`
	NotNull    bool   `gorm:"default:false" mapstructure:"notNull" json:"notNull"`
}

type orm_association struct {
	gorm.Model
	Name                 string    `gorm:"type:varchar(255);not null" mapstructure:"name" json:"name"`
	Type                 string    `gorm:"type:varchar(255);not null" mapstructure:"type" json:"type"`
	TableName            string    `gorm:"type:varchar(255);not null" mapstructure:"tableName" json:"tableName"`
	ColumnName           string    `gorm:"type:varchar(255);not null" mapstructure:"columnName" json:"columnName"`
	ReferencesTable      string    `gorm:"type:varchar(255);not null" mapstructure:"referencesTable" json:"referencesTable"`
	ReferencesField      string    `gorm:"type:varchar(255);not null" mapstructure:"referencesField" json:"referencesField"`
	JoinTable            string    `gorm:"type:varchar(255);not null" mapstructure:"joinTable" json:"joinTable"`
	JoinTableSourceField string    `gorm:"type:varchar(255);not null" mapstructure:"joinTableSourceField" json:"joinTableSourceField"`
	JoinTableTargetField string    `gorm:"type:varchar(255);not null" mapstructure:"joinTableTargetField" json:"joinTableTargetField"`
	OwningModelID        uint      `gorm:"not null" mapstructure:"owningModelID" json:"owningModelID"`
	ReferencesModelID    uint      `gorm:"not null" mapstructure:"referencesModelID" json:"referencesModelID"`
	OwningModel          orm_model `gorm:"foreignKey:OwningModelID" mapstructure:"owningModel" json:"owningModel"`
	ReferencesModel      orm_model `gorm:"foreignKey:ReferencesModelID" mapstructure:"referencesModel" json:"referencesModel"`
}

type response_profile struct {
	gorm.Model
	Name               string       `gorm:"not null" json:"trigger" mapstructure:"trigger"`
	ResponseHTTPStatus int          `gorm:"not null" json:"responseHTTPStatus" mapstructure:"responseHTTPStatus"`
	BodyFormat         pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null" json:"bodyFormat" mapstructure:"bodyFormat"`
}

type internal_tag_group struct {
	gorm.Model
	Name string          `json:"name" mapstructure:"name"`
	Tags []*internal_tag `gorm:"many2many:internal_tag_group_mapping;" json:"tags" mapstructure:"tags"`
}

type internal_tag struct {
	gorm.Model
	Name     string                `gorm:"unique" json:"name" mapstructure:"name"`
	Groups   []*internal_tag_group `gorm:"many2many:internal_tag_group_mapping;" json:"groups" mapstructure:"groups"`
	All      bool                  `json:"all" mapstructure:"all"`
	Reserved bool                  `json:"reserved" mapstructure:"reserved"`
}
