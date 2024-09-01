package postgres

import (
	"gorm.io/gorm"
)

type PostgresBaseRepository struct {
	client *gorm.DB
}

func NewPostgresBaseRepository(client *gorm.DB) *PostgresBaseRepository {
	if client == nil {
		panic("missing postgres client")
	}
	return &PostgresBaseRepository{client: client}
}
