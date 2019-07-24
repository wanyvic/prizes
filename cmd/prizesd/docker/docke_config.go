package dockerapi

import (
	"github.com/docker/docker/client"
)

var (
	DefaultDockerAPIVersion = "1.38"
	cli                     *client.Client
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
