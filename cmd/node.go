package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

func GetNodeInfo(NodeID string) (*swarm.Node, error) {
	node, _, err := dockerapi.CLI.NodeInspectWithRaw(context.Background(), NodeID)
	if err != nil {
		return nil, err
	}
	return &node, nil
}
func RemoveNode(NodeID string, force bool) error {
	err := dockerapi.CLI.NodeRemove(context.Background(), NodeID, types.NodeRemoveOptions{Force: force})
	if err != nil {
		return err
	}
	return nil
}
