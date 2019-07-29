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
	UpdatePrizesServiceOne(service.PrizesService) (bool, error)
	UpdateTaskOne(swarm.Task) (bool, error)
	FindNodeOne(NodeID string) (*swarm.Node, error)
	FindServiceOne(serviceID string) (*swarm.Service, error)
	FindPrizesServiceOne(serviceID string) (*service.PrizesService, error)
	FindTaskList(serviceID string) (*[]swarm.Task, error)
	// UpdateServiceOrderOne(prizestypes.ServiceOrder) (bool, error)
	// FindServiceOrderOne(orderID string) (*prizestypes.ServiceOrder, error)
}

var (
	DBimplement DataBase
)

func init() {
	DBimplement = &mongodb.MongDBClient{URI: mongodb.MongoDBDefaultURI, DataBase: mongodb.MongoDBDefaultDataBase}
}
