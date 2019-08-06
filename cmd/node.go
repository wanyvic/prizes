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
func GetNodeList() (*prizestypes.NodeListStatistics, error) {
	var nodeListStatistics prizestypes.NodeListStatistics
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return nil, err
	}

	swarminfo, err := cli.SwarmInspect(context.Background())
	if err != nil {
		return nil, err
	}

	nodeListStatistics.WorkerToken = swarminfo.JoinTokens.Worker
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
		if node.Spec.Role == swarm.NodeRoleManager {
			if nodeListStatistics.WorkerToken != "" {
				nodeListStatistics.WorkerToken += " " + node.ManagerStatus.Addr
			}
			continue
		}
		nodeListStatistics.TotalCount++
		if node.Status.State == swarm.NodeStateReady {
			nodeListStatistics.AvailabilityCount++
		}
		if status[node.ID] == 0 {
			nodeListStatistics.List = append(nodeListStatistics.List, parseNodeInfo(&node, false))
			nodeListStatistics.UsableCount++
		} else {
			nodeListStatistics.List = append(nodeListStatistics.List, parseNodeInfo(&node, true))
		}
	}
	return &nodeListStatistics, nil
}
func removeNode(NodeID string, force bool) error {
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
func parseNodeInfo(node *swarm.Node, OnWorking bool) (ref prizestypes.NodeInfo) {
	ref.NodeID = node.ID
	ref.Labels = node.Description.Engine.Labels
	ref.NodeState = string(node.Status.State)
	ref.ReachAddress = node.Description.Engine.Labels["REVENUE_ADDRESS"]
	ref.Hardware.CPUType = node.Description.Engine.Labels["CPUNAME"]
	ref.Hardware.CPUThread, _ = strconv.ParseInt(node.Description.Engine.Labels["CPUCOUNT"], 10, 64)
	ref.Hardware.MemoryType = node.Description.Engine.Labels["MEMNAME"]
	ref.Hardware.MemoryCount, _ = strconv.ParseInt(node.Description.Engine.Labels["MEMCOUNT"], 10, 64)
	ref.Hardware.GPUType = node.Description.Engine.Labels["GPUNAME"]
	ref.Hardware.GPUCount, _ = strconv.ParseInt(node.Description.Engine.Labels["GPUCOUNT"], 10, 64)
	ref.Hardware.PersistentStore = node.Description.Engine.Labels["NFSIP"]
	ref.OnWorking = OnWorking
	return ref
}
