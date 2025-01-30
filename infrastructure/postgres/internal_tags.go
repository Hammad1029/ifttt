package postgres

import (
	configuration "ifttt/manager/domain/configuration"

	"gorm.io/gorm"
)

type PostgresInternalTagRepository struct {
	*PostgresBaseRepository
}

func NewPostgresInternalTagRepository(base *PostgresBaseRepository) *PostgresInternalTagRepository {
	return &PostgresInternalTagRepository{PostgresBaseRepository: base}
}

func (p *PostgresInternalTagRepository) AddGroup(dGroup *configuration.InternalTagGroup) error {
	var pgGroup internal_tag_group
	pgGroup.fromDomain(dGroup)
	if err := p.client.Create(&pgGroup).Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresInternalTagRepository) GetAllGroups() (*[]configuration.InternalTagGroup, error) {
	var pgGroups []internal_tag_group
	if err := p.client.Preload("Tags").Find(&pgGroups).Error; err != nil {
		return nil, err
	}

	var dGroups []configuration.InternalTagGroup
	for _, g := range pgGroups {
		dGroups = append(dGroups, *g.toDomain())
	}
	return &dGroups, nil
}

func (p *PostgresInternalTagRepository) GetGroupByName(name string) (*configuration.InternalTagGroup, error) {
	var pgGroup internal_tag_group
	if err := p.client.Preload("Tags").Where("name = ?", name).First(&pgGroup).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	dGroup := pgGroup.toDomain()
	return dGroup, nil
}

func (p *PostgresInternalTagRepository) Add(pTag *configuration.InternalTag) error {
	var pgPTag internal_tag
	pgPTag.fromDomain(pTag)
	if err := p.client.Create(&pgPTag).Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresInternalTagRepository) GetByIDOrName(id uint, name string) (*configuration.InternalTag, error) {
	var pgPTag internal_tag
	if err := p.client.Preload("Groups").
		Where("id = ? or name = ?", id, name).First(&pgPTag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	dPTag := pgPTag.toDomain()
	return dPTag, nil
}

func (p *PostgresInternalTagRepository) GetAll() (*[]configuration.InternalTag, error) {
	var pgPTag []internal_tag
	if err := p.client.Preload("Groups").Find(&pgPTag).Error; err != nil {
		return nil, err
	}

	dPTags := make([]configuration.InternalTag, len(pgPTag))
	for idx, p := range pgPTag {
		dPTags[idx] = *p.toDomain()
	}
	return &dPTags, nil
}
