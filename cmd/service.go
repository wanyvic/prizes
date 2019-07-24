package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
)

func CreateService(serviceSpec swarm.ServiceSpec, options types.ServiceCreateOptions) (*types.ServiceCreateResponse, error) {
	logrus.Info("request CreateService")
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return nil, err
	}
	response, err := cli.ServiceCreate(context.Background(), serviceSpec, options)
	if err != nil {
		return nil, err
	}
	logrus.Info(fmt.Sprintf("CreateService completed: ID: %s ,Warning: %s", response.ID, response.Warnings))
	return &response, nil
}
func UpdateService(serviceID string, serviceSpec swarm.ServiceSpec, options types.ServiceUpdateOptions) (*types.ServiceUpdateResponse, error) {
	logrus.Info("request UpdateService: ", serviceID)
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return nil, err
	}
	service, _, err := cli.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return nil, err
	}
	response, err := cli.ServiceUpdate(context.Background(), service.ID, service.Version, serviceSpec, types.ServiceUpdateOptions{})
	if err != nil {
		return nil, err
	}
	logrus.Info(fmt.Sprintf("CreateService completed: ID: %s ,Warning: %s", serviceID, response.Warnings))
	return &response, nil
}

func RemoveService(serviceID string) error {
	logrus.Info("request RemoveService: ", serviceID)
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	err = cli.ServiceRemove(context.Background(), serviceID)
	if err != nil {
		return err
	}
	err = refresh.RefreshStopService(serviceID)
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("RemoveService completed: ID: %s", serviceID))
	return nil
}
func ServiceInfo(serviceID string) (*swarm.Service, error) {
	logrus.Info("request ServiceInfo: ", serviceID)
	service, err := db.DBimplement.FindServiceOne(serviceID)
	if err != nil {
		return service, err
	}
	return service, nil
}
func TasksInfo(serviceID string) (*[]swarm.Task, error) {
	logrus.Info("request TasksInfo: ", serviceID)
	taskList, err := db.DBimplement.FindTaskList(serviceID)
	if err != nil {
		return nil, err
	}
	return taskList, nil
}
