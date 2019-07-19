package mongodb

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func Test_Connect(t *testing.T) {
	mgo := MongDBClient{URI: MongoDBDefaultURI, DataBase: MongoDBDefaultDataBase}

	cli, err := client.NewClient(client.DefaultDockerHost, "1.38", nil, map[string]string{"Content-type": "application/x-tar"})
	if err != nil {
		t.Error(err)
	}
	servicelist, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		t.Error(err)
	}
	for _, service := range servicelist {
		service.Spec.Annotations.Name = "1111111111111"
		if _, err := mgo.UpdateServiceOne(service); err != nil {
			t.Error(err)
		}

	}
}
func Test_op(t *testing.T) {
	fmt.Println("test_op")

	mgo := MongDBClient{URI: MongoDBDefaultURI, DataBase: MongoDBDefaultDataBase}
	service, err := mgo.FindServiceOne("qrjvzqos4y49nvwp6akswuck6")

	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("service: %s\n", service.ID)
	tasklist, err := mgo.FindTaskList(service.ID)
	if err != nil {
		t.Error(err)
		return
	}
	var td time.Duration
	for _, task := range *tasklist {
		fmt.Printf("\ttask: %s %s %s\n", task.ID, task.DesiredState, task.Status.Timestamp.Sub(task.CreatedAt))
		td += task.Status.Timestamp.Sub(task.CreatedAt)
	}
	fmt.Printf("service totsal avtime %s\n", td)
}
