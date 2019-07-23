package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
)

func CreateService(serviceSpec swarm.ServiceSpec, options types.ServiceCreateOptions) (*types.ServiceCreateResponse, error) {
	response, err := dockerapi.CLI.ServiceCreate(context.Background(), serviceSpec, options)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
func UpdateService(serviceID string, serviceSpec swarm.ServiceSpec, options types.ServiceUpdateOptions) (*types.ServiceUpdateResponse, error) {
	service, _, err := dockerapi.CLI.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return nil, err
	}
	response, err := dockerapi.CLI.ServiceUpdate(context.Background(), service.ID, service.Version, serviceSpec, types.ServiceUpdateOptions{})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func RemoveService(serviceID string) error {
	logrus.Info("RemoveService: ", serviceID)
	err := dockerapi.CLI.ServiceRemove(context.Background(), serviceID)
	if err != nil {
		return err
	}
	err = refresh.RefreshStopService(serviceID)
	if err != nil {
		return err
	}
	return nil
}
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
