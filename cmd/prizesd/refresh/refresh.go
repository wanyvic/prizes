package refresh

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

func RefreshStopService(serviceID string) error {

	fmt.Println("RefreshStopService")
	taskList, err := db.DBimplement.FindTaskList(serviceID)
	if err != nil {
		return err
	}
	for _, task := range *taskList {
		if task.DesiredState == "running" {
			task.Status.Timestamp = time.Now()
			task.DesiredState = "shutdown"
			if _, err := db.DBimplement.UpdateTaskOne(task); err != nil {
				return err
			}
		}
	}
	return nil
}
func RefreshService(serviceID string) error {
	fmt.Println("RefreshService")
	service, _, err := dockerapi.CLI.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return err
	}
	if _, err := db.DBimplement.UpdateServiceOne(service); err != nil {
		return err
	}
	validNameFilter := filters.NewArgs()
	validNameFilter.Add("service", service.ID)
	taskList, err := dockerapi.CLI.TaskList(context.Background(), types.TaskListOptions{Filters: validNameFilter})
	if err != nil {
		return err
	}
	for _, task := range taskList {
		fmt.Printf("\ttask: %s\n", task.ID)
		if _, err := db.DBimplement.UpdateTaskOne(task); err != nil {
			return err
		}
	}
	return nil
}
func WhileLoop() error {

	fmt.Println("WhileLoop")
	for {
		if err := refreshDockerNode(); err != nil {
			return err
		}
		if err := refreshDockerService(); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		if dockerapi.FExit {
			return nil
		}
	}
}
func refreshDockerService() error {
	fmt.Println("refreshDockerService")
	servicelist, err := dockerapi.CLI.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return err
	}
	for _, service := range servicelist {
		fmt.Printf("\tservice: %s\n", service.ID)
		if _, err := db.DBimplement.UpdateServiceOne(service); err != nil {
			return err
		}
		if err := refreshDockerTaskFromService(service.ID); err != nil {
			return err
		}
	}
	return nil
}
func refreshDockerNode() error {
	fmt.Println("refreshDockerNode")
	nodelist, err := dockerapi.CLI.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return err
	}
	for _, node := range nodelist {
		fmt.Printf("\tnode: %s\n", node.ID)
		if _, err := db.DBimplement.UpdateNodeOne(node); err != nil {
			return err
		}
	}
	return nil
}
func refreshDockerTaskFromService(serviceID string) error {
	fmt.Println("refreshDockerTaskFromService")
	validNameFilter := filters.NewArgs()
	validNameFilter.Add("service", serviceID)
	tasklist, err := dockerapi.CLI.TaskList(context.Background(), types.TaskListOptions{Filters: validNameFilter})
	if err != nil {
		return err
	}
	for _, task := range tasklist {

		fmt.Printf("\t\ttask: %s\n", task.ID)
		if _, err := db.DBimplement.UpdateTaskOne(task); err != nil {
			return err
		}
	}
	return nil
}
