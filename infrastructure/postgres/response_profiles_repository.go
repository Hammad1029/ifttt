package postgres

import (
	configuration "ifttt/manager/domain/configuration"

	"gorm.io/gorm"
)

type PostgresResponseProfileRepository struct {
	*PostgresBaseRepository
}

func NewPostgresResponseProfileRepository(base *PostgresBaseRepository) *PostgresResponseProfileRepository {
	return &PostgresResponseProfileRepository{PostgresBaseRepository: base}
}

func (r *PostgresResponseProfileRepository) AddProfile(p *configuration.ResponseProfile) error {
	var pgProfile response_profile
	if err := pgProfile.fromDomain(p); err != nil {
		return err
	}
	if err := r.client.Create(&pgProfile).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresResponseProfileRepository) GetProfilesByName(name string) (*[]configuration.ResponseProfile, error) {
	var pgProfiles []response_profile
	if err := r.client.
		Where("name = ?", name).
		Find(&pgProfiles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	dProfiles := make([]configuration.ResponseProfile, 0, len(pgProfiles))
	for _, p := range pgProfiles {
		if dP, err := p.toDomain(); err != nil {
			return nil, err
		} else {
			dProfiles = append(dProfiles, *dP)
		}
	}
	return &dProfiles, nil
}

func (r *PostgresResponseProfileRepository) GetAllProfiles() (*[]configuration.ResponseProfile, error) {
	var pgProfiles []response_profile
	if err := r.client.Find(&pgProfiles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	dProfiles := make([]configuration.ResponseProfile, 0, len(pgProfiles))
	for _, p := range pgProfiles {
		if dP, err := p.toDomain(); err != nil {
			return nil, err
		} else {
			dProfiles = append(dProfiles, *dP)
		}
	}
	return &dProfiles, nil
}
