package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

func createService(serviceSpec swarm.ServiceSpec, options types.ServiceCreateOptions) error {
	_, err := dockerapi.CLI.ServiceCreate(context.Background(), serviceSpec, options)
	if err != nil {
		return err
	}
	return nil
}
