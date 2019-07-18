package refreshdata

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/wanyvic/prizes/cmd/db"
)

var (
	fExit                   = false
	DefaultDockerAPIVersion = "1.38"
	cli                     *client.Client
)

func exit() {
	fExit = true
}
func init() {
	var err error
	cli, err = client.NewClient(client.DefaultDockerHost, DefaultDockerAPIVersion, nil, map[string]string{"Content-type": "application/x-tar"})
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		exit()
	}()
}
func WhileLoop() error {

	fmt.Println("WhileLoop")
	for {

		servicelist, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
		if err != nil {
			return err
		}
		fmt.Println("update data")
		for _, service := range servicelist {
			fmt.Printf("service: %s\n", service.ID)
			if _, err := db.MongDBClient.UpdateServiceOne(service); err != nil {
				return err
			}
			validNameFilter := filters.NewArgs()
			validNameFilter.Add("service", service.ID)
			tasklist, err := cli.TaskList(context.Background(), types.TaskListOptions{Filters: validNameFilter})
			if err != nil {
				return err
			}
			for _, task := range tasklist {

				fmt.Printf("\ttask: %s\n", task.ID)
				if _, err := db.MongDBClient.UpdateTaskOne(task); err != nil {
					return err
				}
			}
		}
		time.Sleep(1 * time.Second)
		if fExit {
			return nil
		}
	}
}
