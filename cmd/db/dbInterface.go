package db

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/wanyvic/prizes/api/types/service"
	"github.com/wanyvic/prizes/cmd/db/mongodb"
)

type DataBase interface {
	InsertNodeOne(swarm.Node) (bool, error)
	InsertServiceOne(swarm.Service) (bool, error)
	InsertTaskOne(swarm.Task) (bool, error)
	UpdateNodeOne(swarm.Node) (bool, error)
	UpdateServiceOne(swarm.Service) (bool, error)
	UpdateStateTimeAxisOne(service.ServiceTimeLine) (bool, error)
	UpdatePrizesServiceOne(service.PrizesService) (bool, error)
	UpdateTaskOne(swarm.Task) (bool, error)
	FindNodeOne(NodeID string) (*swarm.Node, error)
	FindServiceOne(serviceID string) (*swarm.Service, error)
	FindPrizesServiceOne(serviceID string) (*service.PrizesService, error)
	FindStateTimeAxisOne(serviceID string) (*service.ServiceTimeLine, error)
	FindTaskList(serviceID string) (*[]swarm.Task, error)
	FindPrizesServiceFromPubkey(pubkey string) (*[]service.PrizesService, error)
}

const (
	DBDefaultDataBase = "docker"
)

var (
	DBimplement = &mongodb.MongDBClient{URI: mongodb.MongoDBDefaultURI, DataBase: DBDefaultDataBase}
)
