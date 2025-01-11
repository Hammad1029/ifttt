package postgres

import (
	"fmt"
	"ifttt/manager/domain/user"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	*PostgresBaseRepository
}

func NewPostgresUserRepository(base *PostgresBaseRepository) *PostgresUserRepository {
	return &PostgresUserRepository{PostgresBaseRepository: base}
}

func (p *PostgresUserRepository) GetUser(
	email string, decodeFunc func(input any) (*user.User, error)) (*user.User, error) {

	var pgUser users
	if err := p.client.Where(&users{Email: email}).First(&pgUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("method *PostgresUserRepository.ValidateCredentials: could not query user: %s", err)
	}

	return decodeFunc(pgUser)
}

func (p *PostgresUserRepository) CreateUser(user user.User) error {
	var pgUser users
	if err := mapstructure.Decode(user, &pgUser); err != nil {
		return fmt.Errorf("method *PostgresUserRepository.CreateUser: could not decode user: %s", err)
	}

	if err := p.client.Create(&pgUser).Error; err != nil {
		return fmt.Errorf("method *PostgresUserRepository.CreateUser: could not create user: %s", err)
	}
	return nil
}

func (p *PostgresUserRepository) GetAllUsers() ([]*user.User, error) {
	var pgUsers []users
	if err := p.client.Find(&pgUsers).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresUserRepository.GetUsers: could not get users: %s", err)
	}

	var users []*user.User
	if err := mapstructure.Decode(pgUsers, &users); err != nil {
		return nil, fmt.Errorf("method *PostgresUserRepository.GetUsers: could not decode users: %s", err)
	}

	return users, nil
}
