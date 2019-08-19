package dockerapi

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var (
	DefaultDockerAPIVersion = "1.38"
	cli                     *client.Client
	ProxyImage              = "massgrid/10.0-ubuntu16.04-proxy:v1.0"
	ProxyName               = "massgrid_proxy"
)

func GetDockerClient() (*client.Client, error) {
	if cli == nil {
		var err error
		cli, err = client.NewClient(client.DefaultDockerHost, DefaultDockerAPIVersion, nil, map[string]string{"Content-type": "application/x-tar"})
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}
func DestoryDockerClient() {
	cli.Close()
	cli = nil
}
func NewProxy() error {
	cli, err := GetDockerClient()
	if err != nil {
		return err
	}
	validNameFilter := filters.NewArgs()
	validNameFilter.Add("name", ProxyName)
	servicelist, err := cli.ServiceList(context.Background(), types.ServiceListOptions{Filters: validNameFilter})
	if err != nil {
		return err
	}
	if len(servicelist) > 0 {
		logrus.Info("massgrid_proxy has been started")
		return nil
	}
	spec := ProxyServiceSpec()
	response := types.ServiceCreateResponse{}
	response, err = cli.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
	if err != nil {
		return err
	}
	logrus.Info("massgrid_proxy create successful ", response.ID)
	return nil
}
func ProxyServiceSpec() swarm.ServiceSpec {
	spec := swarm.ServiceSpec{}
	spec.Name = ProxyName
	spec.Mode.Global = &swarm.GlobalService{}
	spec.TaskTemplate.ContainerSpec.Image = ProxyImage
	spec.TaskTemplate.ContainerSpec.Labels = make(map[string]string)
	spec.TaskTemplate.ContainerSpec.Labels["com.massgrid.type"] = "proxy"

	spec.TaskTemplate.Placement = &swarm.Placement{}
	spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, "node.role == worker")
	mount := mount.Mount{Source: "/var/run/docker.sock", Target: "/var/run/docker.sock"}
	spec.TaskTemplate.ContainerSpec.Mounts = append(spec.TaskTemplate.ContainerSpec.Mounts, mount)
	return spec
}
