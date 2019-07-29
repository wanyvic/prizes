package prizeservice

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	"crypto/rand"
	mathRand "math/rand"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

var (
	DefaultDockerImage = "massgrid/10.0-base-ubuntu16.04"
)

// 创建 服务
// 通过 serviceCreate 配置信息创建服务 返回 PrizesService 和错误信息
func Create(serviceCreate *service.ServiceCreate) (*service.PrizesService, *types.ServiceCreateResponse, error) {
	logrus.Info("PrizesService Create")
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
	// prizeService.DockerSerivce, _, err = cli.ServiceInspectWithRaw(context.Background(), response.ID, types.ServiceInspectOptions{})
	// if err != nil {
	// 	return nil, nil, err
	// }

	serviceCreateOrder(&prizeService)

	logrus.Info(fmt.Sprintf("CreateService completed: ID: %s ,Warning: %s", response.ID, response.Warnings))
	return &prizeService, &response, nil
}

func parseServiceCreateSpec(serviceCreate *service.ServiceCreate) *swarm.ServiceSpec {
	replicas := uint64(1)
	spec := swarm.ServiceSpec{}
	spec.TaskTemplate.ContainerSpec = &swarm.ContainerSpec{}
	if len(serviceCreate.ServiceName) > 10 {
		spec.Name = serviceCreate.ServiceName[:10]
		spec.Name += "_" + CreateRandomString(6)
	} else {
		spec.Name = CreateRandomString(10)
	}

	spec.Labels = make(map[string]string)
	spec.Labels["com.massgird.deletetime"] = time.Now().UTC().Add(time.Duration(float64(serviceCreate.Amount)/float64(serviceCreate.ServicePrice)*3600.0) * time.Second).String()
	spec.Labels["com.massgrid.pubkey"] = serviceCreate.Pubkey
	spec.Labels["com.massgrid.price"] = strconv.FormatInt(serviceCreate.ServicePrice, 10)
	spec.Labels["com.massgrid.payment"] = strconv.FormatInt(serviceCreate.Amount, 10)
	spec.Labels["com.massgrid.cputype"] = serviceCreate.CPUType
	spec.Labels["com.massgrid.cputhread"] = strconv.FormatInt(serviceCreate.CPUThread, 10)
	spec.Labels["com.massgrid.memorytype"] = serviceCreate.MemoryType
	spec.Labels["com.massgrid.memorycount"] = strconv.FormatInt(serviceCreate.MemoryCount, 10)
	spec.Labels["com.massgrid.gputype"] = serviceCreate.GPUType
	spec.Labels["com.massgrid.gpucount"] = strconv.FormatInt(serviceCreate.GPUCount, 10)
	spec.Labels["com.massgrid.outpoint.1."+serviceCreate.OutPoint] = strconv.FormatBool(false)

	spec.Mode.Replicated = &swarm.ReplicatedService{Replicas: &replicas}
	if strings.Contains(serviceCreate.Image, "massgrid/") {
		spec.TaskTemplate.ContainerSpec.Image = serviceCreate.Image
	} else {
		spec.TaskTemplate.ContainerSpec.Image = DefaultDockerImage
	}
	spec.TaskTemplate.ContainerSpec.User = "root"

	limits := swarm.GenericResource{DiscreteResourceSpec: &swarm.DiscreteGenericResource{}}
	limits.DiscreteResourceSpec.Kind = serviceCreate.GPUType
	limits.DiscreteResourceSpec.Value = serviceCreate.GPUCount

	spec.TaskTemplate.Resources = &swarm.ResourceRequirements{Reservations: &swarm.Resources{}}
	spec.TaskTemplate.Resources.Reservations.GenericResources = append(spec.TaskTemplate.Resources.Reservations.GenericResources, limits)
	if serviceCreate.SSHPubkey != "" {
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "N2N_SERVERIP="+GetFreeIp().String())
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "N2N_NETMASK=255.0.0.0")
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "N2N_SNIP="+serviceCreate.MasterNodeN2NAddr)
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "SSH_PUBKEY="+serviceCreate.SSHPubkey)
	}

	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "CPUTYPE="+serviceCreate.CPUType)
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "CPUCOUNT="+strconv.FormatInt(serviceCreate.CPUThread, 10))
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "MEMORYTYPE="+serviceCreate.MemoryType)
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "MEMORYCOUNT="+strconv.FormatInt(serviceCreate.MemoryCount, 10))
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "GPUTYPE="+serviceCreate.GPUType)
	spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, "GPUTYPE="+strconv.FormatInt(serviceCreate.GPUCount, 10))
	for k, v := range serviceCreate.ENV {
		spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, strings.ToUpper(k+"="+v))
	}
	//constraints
	spec.TaskTemplate.Placement = &swarm.Placement{}
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "node.role == worker")
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels.cputype  == "+serviceCreate.CPUType)
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels.cputhread == "+strconv.FormatInt(serviceCreate.CPUThread, 10))
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels.memorytype  == "+serviceCreate.MemoryType)
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels.memorycount == "+strconv.FormatInt(serviceCreate.MemoryCount, 10))
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels.gputype  == "+serviceCreate.GPUType)
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "engine.labels.gpucount == "+strconv.FormatInt(serviceCreate.GPUCount, 10))

	mount := mount.Mount{Source: "/dev/net", Target: "/dev/net", ReadOnly: true}
	spec.TaskTemplate.ContainerSpec.Mounts = append(spec.TaskTemplate.ContainerSpec.Mounts, mount)

	return &spec
}

func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyz1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
func GetFreeIp() net.IP {
	mathRand.Seed(time.Now().UnixNano())
	int1 := mathRand.Intn(254)
	mathRand.Seed(time.Now().UnixNano())
	int2 := mathRand.Intn(255)

	var bytes [4]byte
	bytes[0] = byte((int1) & 0xFF)
	bytes[1] = byte((int2) & 0xFF)
	bytes[2] = byte(0x00)
	bytes[3] = byte(0x0A)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func serviceCreateOrder(p *service.PrizesService) {
	p.State = "running"
	p.CreatedAt = p.DockerSerivce.Meta.CreatedAt
	timeScale := time.Duration(float64(p.CreateSpec.Amount) / float64(p.CreateSpec.ServicePrice) * float64(time.Hour))
	p.DeleteAt = p.CreatedAt.Add(timeScale)
	serviceOrder := order.ServiceOrder{}
	serviceOrder.OutPoint = p.CreateSpec.OutPoint
	serviceOrder.CreatedAt = p.CreatedAt
	serviceOrder.RemoveAt = p.DeleteAt
	serviceOrder.OrderState = order.OrderStatePaying
	serviceOrder.Balance = p.CreateSpec.Amount
	serviceOrder.ServicePrice = p.CreateSpec.ServicePrice
	serviceOrder.LastStatementTime = p.CreatedAt
	serviceOrder.NextStatementTime = p.CreatedAt.Add(order.DefaultStatementOptions.StatementDuration)
	p.Order = append(p.Order, serviceOrder)
}
