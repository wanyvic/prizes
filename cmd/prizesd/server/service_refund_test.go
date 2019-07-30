package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/wanyvic/prizes/cmd/prizesd/server"
)

func Test_service_refund(t *testing.T) {
	s := os.Args[len(os.Args)-1]
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", server.DefaultAddr)
			},
		},
	}
	res, err := httpc.Post("http://unix/ServiceRefund/"+s,
		"application/x-www-form-urlencoded", strings.NewReader(""))
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
