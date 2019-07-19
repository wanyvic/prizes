package cmd

import (
	"context"
	"fmt"

	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
)

func removeService(serviceID string) error {
	fmt.Println("removeService: serviceID : ", serviceID)
	err := dockerapi.CLI.ServiceRemove(context.Background(), serviceID)
	if err != nil {
		return err
	}
	err = refresh.RefreshStopService(serviceID)
	if err != nil {
		return err
	}
	return nil
}
