package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
)

func RemoveService(serviceID string) error {
	logrus.Info("RemoveService: ", serviceID)
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
