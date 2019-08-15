package prizeservice

import (
	"bytes"
	"context"
	"crypto/rand"
	"math/big"
	mathRand "math/rand"
	"net"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/wanyvic/prizes/api/types/service"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

func ServiceInfo(prizeService *service.PrizesService) (*service.ServiceInfo, error) {
	serviceInfo := &service.ServiceInfo{}
	serviceInfo.ServiceID = prizeService.DockerService.ID
	serviceInfo.CreatedAt = prizeService.CreatedAt
	serviceInfo.NextCheckTime = prizeService.NextCheckTime
	serviceInfo.Order = prizeService.Order
	serviceInfo.CreateSpec = prizeService.CreateSpec
	serviceInfo.UpdateSpec = prizeService.UpdateSpec
	serviceInfo.State = prizeService.State

	for i := 0; i < len(serviceInfo.Order); i++ {
		serviceInfo.Order[i].Statement = serviceInfo.Order[i].Statement[0:0]
	}

	if serviceInfo.State == service.ServiceStateRunning {
		cli, err := dockerapi.GetDockerClient()
		if err != nil {
			return nil, err
		}
		validNameFilter := filters.NewArgs()
		validNameFilter.Add("service", prizeService.DockerService.ID)
		validNameFilter.Add("desired-state", string(swarm.TaskStateRunning))
		validNameFilter.Add("desired-state", string(swarm.TaskStateAccepted))
		tasklist, err := cli.TaskList(context.Background(), types.TaskListOptions{Filters: validNameFilter})
		if err != nil {
			return nil, err
		}
		if len(tasklist) > 0 {
			serviceInfo.TaskInfo = &tasklist[0]
		}
	}
	return serviceInfo, nil
}

func CreateRandomNumberString(len int) string {
	var container string
	var str = "1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
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
