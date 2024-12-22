package postgres

import (
	"fmt"

	"gorm.io/gorm"
)

type PostgresBaseRepository struct {
	client *gorm.DB
}

func NewPostgresBaseRepository(client *gorm.DB, migrate bool) *PostgresBaseRepository {
	if client == nil {
		panic("missing postgres client")
	}
	if migrate {
		if err := client.AutoMigrate(
			&users{},
			&crons{}, &apis{},
			&trigger_flows{}, &rules{},
			&orm_model{}, &orm_projection{}, &orm_association{},
		); err != nil {
			panic(fmt.Errorf("could not automigrate gorm:%s", err))
		}
	}
	return &PostgresBaseRepository{client: client}
}
