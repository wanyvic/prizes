package prizeservice

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	prizestypes "github.com/wanyvic/prizes/api/types"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

//Create returns prizeservice, response and error
func Create(serviceCreate *service.ServiceCreate) (*service.PrizesService, *types.ServiceCreateResponse, error) {
	serviceSpec := parseServiceCreateSpec(serviceCreate)
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return nil, nil, err
	}
	response := types.ServiceCreateResponse{}
	response, err = cli.ServiceCreate(context.Background(), *serviceSpec, types.ServiceCreateOptions{})
	if err != nil {
		return nil, nil, err
	}

	prizeService := service.PrizesService{CreateSpec: *serviceCreate}
	prizeService.DockerService, _, err = cli.ServiceInspectWithRaw(context.Background(), response.ID, types.ServiceInspectOptions{})
	if err != nil {
		return nil, nil, err
	}

	serviceCreateOrder(&prizeService, serviceCreate)
	logrus.Info(fmt.Sprintf("CreateService completed: ID: %s ,Warning: %s", response.ID, response.Warnings))
	return &prizeService, &response, nil
}

//parseServiceCreateSpec serviceCreate convert to serviceSpec
func parseServiceCreateSpec(serviceCreate *service.ServiceCreate) *swarm.ServiceSpec {
	serviceCreate.ServiceCreateID = strconv.FormatInt(time.Now().UTC().Unix(), 10) + service.DefaultServiceCreateID + CreateRandomNumberString(8)
	replicas := uint64(1)
	spec := swarm.ServiceSpec{}
	spec.TaskTemplate.ContainerSpec = &swarm.ContainerSpec{}
	pos := len(serviceCreate.ServiceName) - 10
	if pos < 0 {
		spec.Name = serviceCreate.ServiceName
		spec.Name += "_" + CreateRandomString(15-len(serviceCreate.ServiceName))
	} else {
		spec.Name = serviceCreate.ServiceName[:10]
		spec.Name += "_" + CreateRandomString(5)
	}

	// parse service labels
	spec.Labels = make(map[string]string)

	timeScale := time.Duration(float64(serviceCreate.Amount) / float64(serviceCreate.ServicePrice) * float64(time.Hour))
	DeleteAt := time.Now().UTC().Add(timeScale)
	spec.Labels["com.massgird.deletetime"] = DeleteAt.String()
	spec.Labels["com.massgrid.pubkey"] = serviceCreate.Pubkey
	spec.Labels["com.massgrid.price"] = strconv.FormatInt(serviceCreate.ServicePrice, 10)
	spec.Labels["com.massgrid.payment"] = strconv.FormatInt(serviceCreate.Amount, 10)
	spec.Labels["com.massgrid.cputype"] = serviceCreate.Hardware.CPUType
	spec.Labels["com.massgrid.cputhread"] = strconv.FormatInt(serviceCreate.Hardware.CPUThread, 10)
	spec.Labels["com.massgrid.memorytype"] = serviceCreate.Hardware.MemoryType
	spec.Labels["com.massgrid.memorycount"] = strconv.FormatInt(serviceCreate.Hardware.MemoryCount, 10)
	spec.Labels["com.massgrid.gputype"] = serviceCreate.Hardware.GPUType
	spec.Labels["com.massgrid.gpucount"] = strconv.FormatInt(serviceCreate.Hardware.GPUCount, 10)
	spec.Labels["com.massgrid.outpoint.1."+serviceCreate.OutPoint] = strconv.FormatBool(false)
	spec.TaskTemplate.ContainerSpec.Labels = make(map[string]string)
	spec.TaskTemplate.ContainerSpec.Labels["com.massgrid.type"] = "worker"
	//parse service image
	if strings.Contains(serviceCreate.Image, "massgrid/") {
		spec.TaskTemplate.ContainerSpec.Image = serviceCreate.Image
	} else {
		spec.TaskTemplate.ContainerSpec.Image = service.DefaultDockerImage
	}

	// parse task user
	spec.TaskTemplate.ContainerSpec.User = "root"

	//parse service replicas
	spec.Mode.Replicated = &swarm.ReplicatedService{Replicas: &replicas}

	//parse service Resources limits

	limits := swarm.GenericResource{DiscreteResourceSpec: &swarm.DiscreteGenericResource{}}
	limits.DiscreteResourceSpec.Kind = serviceCreate.Hardware.GPUType
	limits.DiscreteResourceSpec.Value = serviceCreate.Hardware.GPUCount

	spec.TaskTemplate.Resources = &swarm.ResourceRequirements{Reservations: &swarm.Resources{}}
	spec.TaskTemplate.Resources.Reservations.GenericResources = append(spec.TaskTemplate.Resources.Reservations.GenericResources, limits)

	spec.TaskTemplate.RestartPolicy = &swarm.RestartPolicy{Condition: swarm.RestartPolicyConditionOnFailure}
	//parse environment
	if serviceCreate.SSHPubkey != "" {
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "N2N_SERVERIP="+GetFreeIp().String())
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "N2N_NETMASK=255.0.0.0")
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "N2N_SNIP="+serviceCreate.MasterNodeN2NAddr)
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "SSH_PUBKEY="+serviceCreate.SSHPubkey)
	}

	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, prizestypes.LabelCPUType+"="+serviceCreate.Hardware.CPUType)
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, prizestypes.LabelCPUThread+"="+strconv.FormatInt(serviceCreate.Hardware.CPUThread, 10))
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, prizestypes.LabelMemoryType+"="+serviceCreate.Hardware.MemoryType)
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, prizestypes.LabelMemoryCount+"="+strconv.FormatInt(serviceCreate.Hardware.MemoryCount, 10))
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, prizestypes.LabelGPUType+"="+serviceCreate.Hardware.GPUType)
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, prizestypes.LabelGPUCount+"="+strconv.FormatInt(serviceCreate.Hardware.GPUCount, 10))
	for k, v := range serviceCreate.ENV {
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, strings.ToUpper(k+"="+v))
	}

	//parse service constraints
	spec.TaskTemplate.Placement = &swarm.Placement{}
	platform := swarm.Platform{Architecture: "amd64", OS: "linux"}
	spec.TaskTemplate.Placement.Platforms = append(spec.TaskTemplate.Placement.Platforms, platform)
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "node.role == worker")
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels."+prizestypes.LabelCPUType+" == "+serviceCreate.Hardware.CPUType)
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels."+prizestypes.LabelCPUThread+" == "+strconv.FormatInt(serviceCreate.Hardware.CPUThread, 10))
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels."+prizestypes.LabelMemoryType+" == "+serviceCreate.Hardware.MemoryType)
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels."+prizestypes.LabelMemoryCount+" == "+strconv.FormatInt(serviceCreate.Hardware.MemoryCount, 10))
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels."+prizestypes.LabelGPUType+" == "+serviceCreate.Hardware.GPUType)
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels."+prizestypes.LabelGPUCount+" == "+strconv.FormatInt(serviceCreate.Hardware.GPUCount, 10))

	//parse mount
	mount := mount.Mount{Source: "/dev/net", Target: "/dev/net", ReadOnly: true}
	spec.TaskTemplate.ContainerSpec.Mounts = append(spec.TaskTemplate.ContainerSpec.Mounts, mount)

	return &spec
}

//serviceCreateOrder create the service order
func serviceCreateOrder(p *service.PrizesService, serviceCreate *service.ServiceCreate) {
	p.State = service.ServiceStateRunning
	p.CreatedAt = p.DockerService.Meta.CreatedAt
	timeScale := time.Duration(float64(serviceCreate.Amount) / float64(serviceCreate.ServicePrice) * float64(time.Hour))
	p.DeleteAt = p.CreatedAt.Add(timeScale)
	serviceOrder := order.ServiceOrder{}
	serviceOrder.OriderID = serviceCreate.ServiceCreateID
	serviceOrder.OutPoint = serviceCreate.OutPoint
	serviceOrder.CreatedAt = p.CreatedAt
	serviceOrder.RemoveAt = p.DeleteAt
	serviceOrder.OrderState = order.OrderStatePaying
	serviceOrder.Drawee = serviceCreate.Drawee
	serviceOrder.Balance = serviceCreate.Amount
	serviceOrder.PayAmount = serviceCreate.Amount
	serviceOrder.ServicePrice = serviceCreate.ServicePrice
	serviceOrder.LastStatementTime = p.CreatedAt
	serviceOrder.MasterNodeFeeRate = serviceCreate.MasterNodeFeeRate
	serviceOrder.MasterNodeFeeAddress = serviceCreate.MasterNodeFeeAddress
	serviceOrder.DevFeeRate = serviceCreate.DevFeeRate
	serviceOrder.DevFeeAddress = serviceCreate.DevFeeAddress
	p.NextCheckTime = p.CreatedAt.Add(StatementDuration)
	if p.NextCheckTime.After(p.DeleteAt) {
		p.NextCheckTime = p.DeleteAt
	}
	p.Order = append(p.Order, serviceOrder)
}
