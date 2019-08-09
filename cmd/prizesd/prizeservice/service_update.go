package prizeservice

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

func Update(prizeService *service.PrizesService, serviceUpdate *service.ServiceUpdate) (*types.ServiceUpdateResponse, error) {
	logrus.Info("PrizesService Create")
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return nil, err
	}
	prizeService.UpdateSpec = append(prizeService.UpdateSpec, *serviceUpdate)
	serviceSpec := parseServiceUpdateSpec(&prizeService.DockerService, serviceUpdate)
	response, err := cli.ServiceUpdate(context.Background(), serviceUpdate.ServiceID, prizeService.DockerService.Meta.Version, *serviceSpec, types.ServiceUpdateOptions{})
	if err != nil {
		return nil, err
	}
	serviceUpdateOrder(prizeService, serviceUpdate)

	logrus.Info(fmt.Sprintf("CreateService completed: ID: %s ,Warning: %s", serviceUpdate.ServiceID, response.Warnings))
	return &response, nil
}
func parseServiceUpdateSpec(dockerservice *swarm.Service, serviceUpdate *service.ServiceUpdate) *swarm.ServiceSpec {
	dockerservice.Spec.Labels["com.massgird.deletetime"] = dockerservice.Meta.CreatedAt.Add(time.Duration(float64(serviceUpdate.Amount)/float64(serviceUpdate.ServicePrice)*3600.0) * time.Second).String()
	num := 1
	for k, _ := range dockerservice.Spec.Labels {
		if strings.Contains(k, "com.massgrid.outpoint") {
			num++
		}
	}
	serviceUpdate.ServiceUpdateID = strconv.FormatInt(time.Now().UTC().Unix(), 10) + service.DefaultServiceUpdateID + CreateRandomNumberString(8)
	dockerservice.Spec.Labels["com.massgrid.outpoint."+strconv.Itoa(num)+"."+serviceUpdate.OutPoint] = strconv.FormatBool(false)
	return &dockerservice.Spec
}
func serviceUpdateOrder(p *service.PrizesService, serviceUpdate *service.ServiceUpdate) {
	serviceOrder := order.ServiceOrder{}
	serviceOrder.OriderID = serviceUpdate.ServiceUpdateID
	serviceOrder.OutPoint = serviceUpdate.OutPoint
	serviceOrder.CreatedAt = p.DeleteAt
	serviceOrder.Drawee = serviceUpdate.Drawee
	timeScale := time.Duration(float64(serviceUpdate.Amount) / float64(serviceUpdate.ServicePrice) * float64(time.Hour))
	p.DeleteAt = p.DeleteAt.Add(timeScale)

	serviceOrder.RemoveAt = p.DeleteAt
	serviceOrder.OrderState = order.OrderStateWaitToPay
	serviceOrder.Balance = serviceUpdate.Amount
	p.Order = append(p.Order, serviceOrder)
}
