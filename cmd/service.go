package cmd

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	prizestypes "github.com/wanyvic/prizes/api/types"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
)

var (
	DefaultDockerImage = "massgrid/10.0-base-ubuntu16.04"
)

func CreateService(serviceCreate prizestypes.ServiceCreate, options types.ServiceCreateOptions) (*types.ServiceCreateResponse, error) {
	
}
func UpdateService(serviceUpdate prizestypes.ServiceUpdate, options types.ServiceUpdateOptions) (*types.ServiceUpdateResponse, error) {
	logrus.Info("request UpdateService: ", serviceUpdate.ServiceID)
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return nil, err
	}
	service, _, err := cli.ServiceInspectWithRaw(context.Background(), serviceUpdate.ServiceID, types.ServiceInspectOptions{})
	if err != nil {
		return nil, err
	}
	serviceSpec := preaseServiceUpdateSpec(&service, &serviceUpdate)
	response, err := cli.ServiceUpdate(context.Background(), service.ID, service.Version, *serviceSpec, types.ServiceUpdateOptions{})
	if err != nil {
		return nil, err
	}

	serviceOrder := preaseUpdateServiceOrder(&serviceUpdate, &response)
	_, err = db.DBimplement.UpdateServiceOrderOne(*serviceOrder)
	if err != nil {
		return nil, err
	}
	logrus.Info(fmt.Sprintf("UpdateService completed: ID: %s ,Warning: %s", serviceUpdate.ServiceID, response.Warnings))
	return &response, nil
}

func ServiceRemove(serviceID string) error {
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
func Service(serviceID string) (*swarm.Service, error) {
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

