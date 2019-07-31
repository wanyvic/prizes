package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/wanyvic/prizes/api/types/service"
	"github.com/wanyvic/prizes/cmd/prizesd/server"
)

func Test_service_create(t *testing.T) {
	serviceCreate := service.ServiceCreate{}
	serviceCreate.ServiceName = "wany"
	serviceCreate.SSHPubkey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsGRtaxyGHIXmgG5dcuj49mlHY2Wg0UYBav+VKjSoESFW0ctRWWEKFWL6I08BSOXimz2HDNFqdyqn8RhCKzohSEDPV9pfmUv6LzV4GJgEFLsuVr4WrSOm9B6a6I6/ckE2iZIIhhdyeCMm6zHBOLohyWNvLIuFy3a5s0nJ7iP1Z/ifLO0Q8FhmAnZtXZX1TqOm7rAXXRAwJ2EFjt8ELuU+ygcdEAl0zjT7QH1JxaXKzDABfT8TPUZkq64QWbgfcNO59mvY7JHOUkZlpy5k6tfFMuNXt0NLtVjRPH0weYKA4iaBSooFloimuJxAQ6cHol4f5RlL0GXacruw5xQVOW+XD wany@WANY"
	serviceCreate.Pubkey = "3f0a6df6b7dd8f265a56f836c3fd166261341e4eeaa6bb021cd754204f761e4b"
	serviceCreate.Amount = 1000000000
	serviceCreate.Image = "massgrid/nginx:latest"
	serviceCreate.ServicePrice = 100000000
	serviceCreate.Drawee = "MQbZKpL53knrkJD13GgYf7hazwfqqGHGDk"
	serviceCreate.BlockHeight = 187291
	serviceCreate.OutPoint = "3f0a6df6b7dd8f265a56f836c3fd166261341e4eeaa6bb021cd754204f761e4b-1"

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", server.DefaultAddr)
			},
		},
	}
	jsonstr, err := json.Marshal(serviceCreate)
	if err != nil {
		t.Error(err)
	}
	res, err := httpc.Post("http://unix/ServiceCreate",
		"application/x-www-form-urlencoded",
		strings.NewReader(string(jsonstr)))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(res.Status)
	for k, v := range res.Header {
		fmt.Println(k, ": ", v)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	strBody := string(body)

	fmt.Println(strBody)
}
