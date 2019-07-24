package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

func GetNodeInfo(NodeID string) (*swarm.Node, error) {
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return nil, err
	}
	node, _, err := cli.NodeInspectWithRaw(context.Background(), NodeID)
	if err != nil {
		return nil, err
	}
	return &node, nil
}
func RemoveNode(NodeID string, force bool) error {
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	err = cli.NodeRemove(context.Background(), NodeID, types.NodeRemoveOptions{Force: force})
	if err != nil {
		return err
	}
	return nil
}
