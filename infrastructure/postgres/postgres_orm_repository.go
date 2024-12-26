package postgres

import (
	"ifttt/manager/domain/orm_schema"

	"gorm.io/gorm"
)

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

func (o *PostgresOrmRepository) GetModelByIdOrName(id uint, name string) (*orm_schema.Model, error) {
	var pgModel orm_model
	if err := o.client.
		Preload("Projections").Preload("OwningAssociations").Preload("ReferencedAssociations").
		First(&pgModel, "name = ? or id = ?", name, id).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
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
		First(&pgAssociation, "name = ?", name).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if dAssociation, err := pgAssociation.toDomain(); err != nil {
		return nil, err
	} else {
		return dAssociation, nil
	}
}

func (o *PostgresOrmRepository) GetAllAssociations() (map[string]*orm_schema.ModelAssociation, error) {
	var pgAssociation []orm_association
	if err := o.client.Find(&pgAssociation).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	mapped := make(map[string]*orm_schema.ModelAssociation)
	for _, a := range pgAssociation {
		if dAssociation, err := a.toDomain(); err != nil {
			return nil, err
		} else {
			mapped[a.Name] = dAssociation
		}
	}
	return mapped, nil
}

func (o *PostgresOrmRepository) GetAllModels() (map[string]*orm_schema.Model, error) {
	var pgModels []orm_model
	if err := o.client.
		Preload("Projections").Preload("OwningAssociations").Preload("ReferencedAssociations").
		Preload("OwningAssociations.ReferencesModel").Preload("ReferencedAssociations.OwningModel").
		Find(&pgModels).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	mapped := make(map[string]*orm_schema.Model)
	for _, a := range pgModels {
		if dModel, err := a.toDomain(); err != nil {
			return nil, err
		} else {
			mapped[a.Name] = dModel
		}
	}
	return mapped, nil
}
