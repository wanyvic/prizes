package dockerapi

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/client"
)

var (
	FExit                   = false
	defaultDockerAPIVersion = "1.38"
	CLI                     *client.Client
)

func init() {
	var err error
	CLI, err = client.NewClient(client.DefaultDockerHost, defaultDockerAPIVersion, nil, map[string]string{"Content-type": "application/x-tar"})
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		FExit = true
	}()
}
