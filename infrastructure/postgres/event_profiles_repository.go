package postgres

import (
	eventprofiles "ifttt/manager/domain/event_profiles"

	"gorm.io/gorm"
)

type PostgresEventProfileRepository struct {
	*PostgresBaseRepository
}

func NewPostgresEventProfileRepository(base *PostgresBaseRepository) *PostgresEventProfileRepository {
	return &PostgresEventProfileRepository{PostgresBaseRepository: base}
}

func (r *PostgresEventProfileRepository) AddProfile(p *eventprofiles.Profile, parent uint) error {
	var pgProfile event_profile
	if err := pgProfile.fromDomain(p); err != nil {
		return err
	}
	if parent != 0 {
		pgProfile.ParentID = &parent
	}
	if err := r.client.Create(&pgProfile).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresEventProfileRepository) GetProfilesByInternalAndTrigger(internal bool, trigger string) (*[]eventprofiles.Profile, error) {
	var pgProfiles []event_profile
	if err := r.client.
		Where("internal = ? AND trigger = ?", internal, trigger).
		Find(&pgProfiles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	dProfiles := make([]eventprofiles.Profile, 0, len(pgProfiles))
	for _, p := range pgProfiles {
		if dP, err := p.toDomain(); err != nil {
			return nil, err
		} else {
			dProfiles = append(dProfiles, *dP)
		}
	}
	return &dProfiles, nil
}

func (r *PostgresEventProfileRepository) GetAllInternalProfiles() (*[]eventprofiles.Profile, error) {
	var pgProfiles []event_profile
	if err := r.client.
		Preload("MappedProfiles", "parent_id IS NOT NULL").
		Where("parent_id IS NULL").Find(&pgProfiles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	dProfiles := make([]eventprofiles.Profile, 0, len(pgProfiles))
	for _, p := range pgProfiles {
		if dP, err := p.toDomain(); err != nil {
			return nil, err
		} else {
			dProfiles = append(dProfiles, *dP)
		}
	}
	return &dProfiles, nil
}
