package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/wanyvic/prizes/cmd/prizesd/server"
)

func Test_service_getservicesfrompubkey(t *testing.T) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", server.DefaultAddr)
			},
		},
	}
	s := "3f0a6df6b7dd8f265a56f836c3fd166261341e4eeaa6bb021cd754204f761e4b"
	res, err := httpc.Post("http://unix/getservicesfrompubkey/"+s,
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
