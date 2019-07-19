package cmd

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/wanyvic/prizes/cmd/db"
)

func ServiceInfo(serviceID string) (*swarm.Service, error) {
	service, err := db.DBimplement.FindServiceOne(serviceID)
	if err != nil {
		return service, err
	}
	return service, nil
}
func TasksInfo(serviceID string) (*[]swarm.Task, error) {
	taskList, err := db.DBimplement.FindTaskList(serviceID)
	if err != nil {
		return nil, err
	}
	return taskList, nil
}
