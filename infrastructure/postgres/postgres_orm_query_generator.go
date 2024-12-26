package postgres

import (
	"ifttt/manager/domain/orm_schema"
	"ifttt/manager/domain/resolvable"
)

type PostgresOrmQueryGeneratorRepository struct {
	*PostgresBaseRepository
}

func NewPostgresOrmQueryGeneratorRepository(base *PostgresBaseRepository) *PostgresOrmQueryGeneratorRepository {
	return &PostgresOrmQueryGeneratorRepository{PostgresBaseRepository: base}
}

func (p *PostgresOrmQueryGeneratorRepository) Generate(
	r *resolvable.OrmResolvable, models map[string]*orm_schema.Model) (string, error) {
	return "", nil
}
