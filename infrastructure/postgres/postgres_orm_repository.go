package postgres

import "ifttt/manager/domain/orm_schema"

type PostgresOrmRepository struct {
	*PostgresBaseRepository
}

func NewPostgresOrmRepository(base *PostgresBaseRepository) *PostgresOrmRepository {
	return &PostgresOrmRepository{PostgresBaseRepository: base}
}

func (o *PostgresOrmRepository) CreateModel(model *orm_schema.Model) error {
	var pgModel orm_model
	if err := pgModel.fromDomain(model); err != nil {
		return err
	}
	if err := o.client.Create(&pgModel).Error; err != nil {
		return err
	}
	return nil
}

func (o *PostgresOrmRepository) GetModelByName(name string) (*orm_schema.Model, error) {
	var pgModel orm_model
	if err := o.client.
		Preload("Projections").Preload("Constraints").
		First(&pgModel, "name = ?", name).Error; err != nil {
		return nil, err
	}
	if dModel, err := pgModel.toDomain(); err != nil {
		return nil, err
	} else {
		return dModel, nil
	}
}

func (o *PostgresOrmRepository) CreateAssociation(association *orm_schema.ModelAssociation) error {
	var pgAssociation orm_association
	if err := pgAssociation.fromDomain(association); err != nil {
		return err
	}
	if err := o.client.Create(&pgAssociation).Error; err != nil {
		return err
	}
	return nil
}

func (o *PostgresOrmRepository) GetAssociationByName(name string) (*orm_schema.ModelAssociation, error) {
	var pgAssociation orm_association
	if err := o.client.
		First(&pgAssociation, "name = ?", name).Error; err != nil {
		return nil, err
	}
	if dAssociation, err := pgAssociation.toDomain(); err != nil {
		return nil, err
	} else {
		return dAssociation, nil
	}
}
