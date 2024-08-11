package scylla

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

type ScyllaBaseRepository struct {
	session *gocqlx.Session
	cluster *gocql.ClusterConfig
}

func NewScyllaBaseRepository(session *gocqlx.Session, cluster *gocql.ClusterConfig) *ScyllaBaseRepository {
	if session == nil {
		panic("missing scylla session")
	}
	if cluster == nil {
		panic("missing scylla cluster")
	}
	return &ScyllaBaseRepository{
		session: session,
		cluster: cluster,
	}
}
