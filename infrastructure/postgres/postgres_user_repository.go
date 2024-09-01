package postgres

import (
	"fmt"
	"ifttt/manager/domain/user"

	"gorm.io/gorm"
)

type postgresUser struct {
	gorm.Model
	Email    string `gorm:"type:varchar(50);unique" mapstructure:"username"`
	Password string `gorm:"type:varchar(50)" mapstructure:"username"`
}

func (p postgresUser) TableName() string {
	return "users"
}

type PostgresUserRepository struct {
	*PostgresBaseRepository
}

func NewPostgresUserRepository(base *PostgresBaseRepository) *PostgresUserRepository {
	return &PostgresUserRepository{PostgresBaseRepository: base}
}

func (p *PostgresUserRepository) GetUser(
	email string, decodeFunc func(input any) (*user.User, error)) (*user.User, error) {

	var pgUser postgresUser
	if err := p.client.Where(&postgresUser{Email: email}).First(&pgUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, fmt.Errorf("method *PostgresUserRepository.ValidateCredentials: could not query user: %s", err)
		}
	}

	return decodeFunc(pgUser)
}
