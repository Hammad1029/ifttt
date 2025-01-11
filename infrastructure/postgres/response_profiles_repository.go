package postgres

import (
	responseprofiles "ifttt/manager/domain/response_profiles"

	"gorm.io/gorm"
)

type PostgresResponseProfileRepository struct {
	*PostgresBaseRepository
}

func NewPostgresResponseProfileRepository(base *PostgresBaseRepository) *PostgresResponseProfileRepository {
	return &PostgresResponseProfileRepository{PostgresBaseRepository: base}
}

func (r *PostgresResponseProfileRepository) AddProfile(p *responseprofiles.Profile, parent uint) error {
	var pgProfile response_profile
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

func (r *PostgresResponseProfileRepository) GetProfilesByInternalAndCode(internal bool, code string) (*[]responseprofiles.Profile, error) {
	var pgProfiles []response_profile
	if err := r.client.
		Preload("Code").Preload("Description").Preload("Data").
		Where("internal = ? AND mapped_code = ?", internal, code).
		Find(&pgProfiles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	dProfiles := make([]responseprofiles.Profile, 0, len(pgProfiles))
	for _, p := range pgProfiles {
		if dP, err := p.toDomain(); err != nil {
			return nil, err
		} else {
			dProfiles = append(dProfiles, *dP)
		}
	}
	return &dProfiles, nil
}

func (r *PostgresResponseProfileRepository) GetAllInternalProfiles() (*[]responseprofiles.Profile, error) {
	var pgProfiles []response_profile
	r.client.
		Preload("Code").Preload("Description").Preload("Data").Preload("Errors").
		Preload("MappedProfiles", "parent_id IS NOT NULL").
		Preload("MappedProfiles.Code").Preload("MappedProfiles.Description").
		Preload("MappedProfiles.Data").Preload("MappedProfiles.Errors").
		Where("parent_id IS NULL").
		Find(&pgProfiles)

	dProfiles := make([]responseprofiles.Profile, 0, len(pgProfiles))
	for _, p := range pgProfiles {
		if dP, err := p.toDomain(); err != nil {
			return nil, err
		} else {
			dProfiles = append(dProfiles, *dP)
		}
	}
	return &dProfiles, nil
}
