package scylla

import (
	"fmt"
	"ifttt/manager/domain/user"

	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/scylladb/gocqlx/v3/table"
)

type scyllaUsers struct {
	Email    string `cql:"email" mapstructure:"email"`
	Password string `cql:"password" mapstructure:"password"`
}

var scyllaUsersMetadata = table.Metadata{
	Name:    "users",
	Columns: []string{"email", "password"},
	PartKey: []string{"group"},
	SortKey: []string{"name"},
}

var scyllaUsersTable *table.Table

type ScyllaUsersRepository struct {
	ScyllaBaseRepository
}

func NewScyllaUsersRepository(base ScyllaBaseRepository) *ScyllaUsersRepository {
	return &ScyllaUsersRepository{ScyllaBaseRepository: base}
}

func (s *ScyllaUsersRepository) getTable() *table.Table {
	if scyllaUsersTable == nil {
		scyllaUsersTable = table.New(scyllaUsersMetadata)
	}
	return scyllaUsersTable
}

func (s *ScyllaUsersRepository) GetUser(
	email string, decodeFunc func(input any) (*user.User, error)) (*user.User, error) {

	var scyllaUser scyllaUsers
	stmt, names := s.getTable().SelectBuilder().Where(qb.Eq("email")).Limit(1).ToCql()
	if err := s.session.Query(stmt, names).BindStruct(scyllaUsers{Email: email}).SelectRelease(&scyllaUser); err != nil {
		return nil, fmt.Errorf("method *ScyllaUsersRepository.ValidateCredentials: could not query user: %s", err)
	}

	return decodeFunc(scyllaUser)
}
