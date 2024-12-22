package infrastructure

import (
	"gorm.io/gorm"
)

type MySqlBaseRepository struct {
	client *gorm.DB
}

func NewMySqlBaseRepository(client *gorm.DB) *MySqlBaseRepository {
	if client == nil {
		panic("missing mysql client")
	}
	return &MySqlBaseRepository{client: client}
}
