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
	serviceCreate.Amount = 3000
	serviceCreate.Image = "massgrid/nginx:latest"
	serviceCreate.ServicePrice = 120000
	serviceCreate.Drawee = "payer"
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
