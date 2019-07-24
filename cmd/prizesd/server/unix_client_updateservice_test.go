package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
)

var serviceSpeca = `{
		"Name": "unruffled_austisn",
		"Labels": {
			"foo": "bar"
		},
		"TaskTemplate": {
			"ContainerSpec": {
      			"Image": "nginx:latest"
			},
			"ForceUpdate": 0,
			"Runtime": "container"
		},
		"Mode": {
			"Replicated": {
				"Replicas": 1
			}
		}
	}`

// func Test_Unix_CreateService_t(t *testing.T) {
// 	httpc := http.Client{
// 		Transport: &http.Transport{
// 			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
// 				return net.Dial("unix", DefaultAddr)
// 			},
// 		},
// 	}
// 	res, err := httpc.Post("http://unix/CreateService",
// 		"application/x-www-form-urlencoded",
// 		strings.NewReader(serviceSpeca))
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}

// 	fmt.Println(res.Status)
// 	for k, v := range res.Header {
// 		fmt.Println(k, ": ", v)
// 	}
// 	defer res.Body.Close()
// 	body, err := ioutil.ReadAll(res.Body)
// 	strBody := string(body)

// 	fmt.Println(strBody)
// }
func Test_Unix_UpdateService_t(t *testing.T) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", DefaultAddr)
			},
		},
	}

	serviceID := "lxfnhu8hidtm56hid4wn2fisj"
	res, err := httpc.Post("http://unix/UpdateService/"+serviceID,
		"application/x-www-form-urlencoded",
		strings.NewReader(serviceSpeca))
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
