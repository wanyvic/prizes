package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/prizeservice"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh/calculagraph"
)

func ServiceCreate(serviceCreate *service.ServiceCreate) (*types.ServiceCreateResponse, error) {
	logrus.Info("request ServiceCreate")
	prizeService, response, err := prizeservice.Create(serviceCreate)
	if err != nil {
		return nil, err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return nil, err
	}
	calculagraph.Push(prizeService.DockerSerivce.ID, prizeService.Order[0].NextStatementTime)
	return response, nil
}
func ServiceUpdate(serviceUpdate *service.ServiceUpdate, options types.ServiceUpdateOptions) (*types.ServiceUpdateResponse, error) {
	logrus.Info("request ServiceUpdate: ", serviceUpdate.ServiceID)

	prizeService, err := db.DBimplement.FindPrizesServiceOne(serviceUpdate.ServiceID)
	if err != nil {
		return nil, err
	}
	response, err := prizeservice.Update(prizeService, serviceUpdate)
	if err != nil {
		return nil, err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return nil, err
	}
	calculagraph.ChangeServiceRemoveTime(prizeService.DockerSerivce.ID, prizeService.DeleteAt)
	return response, nil
}
func serviceSate() {

}
func ServiceStatement(ServiceID string, statementAt time.Time) error {
	logrus.Info("request ServiceStatement: ", ServiceID)
	prizeService, err := db.DBimplement.FindPrizesServiceOne(ServiceID)
	if err != nil {
		return err
	}
	serviceStatistics, err := ServiceState(ServiceID)
	if err != nil {
		return err
	}
	_, err = prizeservice.Statement(prizeService, serviceStatistics, statementAt, order.DefaultStatementOptions)
	if err != nil {
		return err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return err
	}
	calculagraph.ChangeServiceRemoveTime(prizeService.DockerSerivce.ID, prizeService.DeleteAt)
	return nil
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
