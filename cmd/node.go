package cmd

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	prizestypes "github.com/wanyvic/prizes/api/types"
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
func GetNodeList() (*prizestypes.NodeListStatistics, error) {
	var nodeListStatistics prizestypes.NodeListStatistics
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return nil, err
	}
	nodelist, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return nil, err
	}
	validNameFilter := filters.NewArgs()
	validNameFilter.Add("desired-state", "running")
	tasklist, err := cli.TaskList(context.Background(), types.TaskListOptions{Filters: validNameFilter})
	if err != nil {
		return nil, err
	}
	status := make(map[string]int)

	for _, task := range tasklist {
		if _, ok := status["com.massgird.work"]; ok {
			status[task.NodeID]++
		}
	}
	for _, node := range nodelist {
		if node.Spec.Role == "manager" {
			continue
		}
		nodeListStatistics.TotalCount++
		if node.Status.State == "ready" {
			nodeListStatistics.AvailabilityCount++
		}
		if status[node.ID] == 0 {
			nodeListStatistics.List = append(nodeListStatistics.List, preaseNodeInfo(&node, false))
			nodeListStatistics.UsableCount++
		} else {
			nodeListStatistics.List = append(nodeListStatistics.List, preaseNodeInfo(&node, true))
		}
	}
	return &nodeListStatistics, nil
}
func preaseNodeInfo(node *swarm.Node, OnWorking bool) (ref prizestypes.NodeInfo) {
	ref.NodeID = node.ID
	ref.Labels = node.Description.Engine.Labels
	ref.NodeState = string(node.Status.State)
	ref.ReachAddress = node.Description.Engine.Labels["REVENUE_ADDRESS"]
	ref.CPUType = node.Description.Engine.Labels["CPUNAME"]
	ref.CPUThread, _ = strconv.Atoi(node.Description.Engine.Labels["CPUCOUNT"])
	ref.MemoryType = node.Description.Engine.Labels["MEMNAME"]
	ref.MemoryCount, _ = strconv.Atoi(node.Description.Engine.Labels["MEMCOUNT"])
	ref.GPUType = node.Description.Engine.Labels["GPUNAME"]
	ref.GPUCount, _ = strconv.Atoi(node.Description.Engine.Labels["GPUCOUNT"])
	ref.PersistentStore = node.Description.Engine.Labels["NFSIP"]
	ref.OnWorking = OnWorking
	return ref
}
