package db

import (
	"github.com/docker/docker/api/types/swarm"
)

type dbOperator interface {
	InsertServiceOne(swarm.Service) (bool, error)
	UpdateServiceOne(swarm.Service) (bool, error)
	InsertTaskOne(swarm.Task) (bool, error)
	UpdateTaskOne(swarm.Task) (bool, error)
}
