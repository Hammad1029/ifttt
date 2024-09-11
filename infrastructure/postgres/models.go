package postgres

import "gorm.io/gorm"

type postgresAPI struct {
	gorm.Model
	Group       string         `gorm:"type:varchar(50)" mapstructure:"group"`
	Name        string         `gorm:"type:varchar(50)" mapstructure:"name"`
	Method      string         `gorm:"type:varchar(10)" mapstructure:"method"`
	Type        string         `gorm:"type:varchar(10)" mapstructure:"type"`
	Path        string         `gorm:"type:varchar(50)" mapstructure:"path"`
	Description string         `gorm:"type:varchar(255)" mapstructure:"description"`
	Request     map[string]any `gorm:"type:jsonb;default:'{}';not null" mapstructure:"request"`
	Dumping     map[string]any `gorm:"type:jsonb;default:'{}';not null" mapstructure:"dumping"`
	StartRules  []string       `gorm:"type:varchar(50)[];default:'[]';not null" mapstructure:"rules"`
	Rules       map[string]any `gorm:"type:jsonb;default:'{}';not null" mapstructure:"startRules"`
}

func (p postgresAPI) TableName() string {
	return "apis"
}

type postgresUser struct {
	gorm.Model
	Email    string `gorm:"type:varchar(50);unique" mapstructure:"email"`
	Password string `gorm:"type:varchar(255)" mapstructure:"password"`
}

func (p postgresUser) TableName() string {
	return "users"
}
