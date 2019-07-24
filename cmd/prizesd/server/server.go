// +build linux freebsd openbsd darwin

package server // import "github.com/wanyvic/prizes/cmd/prizesd/server"

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

const DefaultProto = "unix"
const DefaultAddr = "/var/run/prizes.sock"
const DefaultHTTPHost = "localhost"
const DefaultHTTPPort = 9333 // Default HTTP Port
var DefaultTCPHost = fmt.Sprintf("%s:%d", DefaultHTTPHost, DefaultHTTPPort)

type ServerOpts struct {
	Proto string
	Addr  string
}
type Server struct {
	proto  string
	addr   string
	server *http.Server
}

func (u *Server) setUnixDefaultOpts() {
	u.proto = DefaultProto
	u.addr = DefaultAddr
}
func (u *Server) setTcpDefaultOpts() {
	u.proto = "tcp"
	u.addr = DefaultTCPHost
}
func NewServer(proto string) (*Server, error) {
	var server Server
	if proto == "unix" {
		server.setUnixDefaultOpts()
	} else if proto == "tcp" {
		server.setTcpDefaultOpts()
	} else {
		return nil, errors.New("undefine proto")
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.Stop()
	}()
	return &server, nil
}
func NewServerWithOpts(opts ServerOpts) (*Server, error) {
	var server Server
	server.proto = opts.Proto
	server.addr = opts.Addr
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.Stop()
	}()
	return &server, nil
}
func (u *Server) Start() {
	logrus.Info(fmt.Sprintf("%s://%s Server Starting", u.proto, u.addr))
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	u.server = &http.Server{
		Handler: mux,
	}
	if strings.Contains(u.proto, "unix") {
		os.Remove(u.addr)
		unixListener, err := net.Listen("unix", u.addr)
		if err != nil {
			logrus.Error(fmt.Sprintf("%s://%s Server start error: %s ", u.proto, u.addr, err))
		}
		u.server.Serve(unixListener)
	} else if strings.Contains(u.proto, "tcp") {
		unixListener, err := net.Listen("tcp", u.addr)
		if err != nil {
			logrus.Error(fmt.Sprintf("%s://%s Server start error: %s ", u.proto, u.addr, err))
		}
		u.server.Serve(unixListener)
	} else {
		logrus.Error(fmt.Sprintf("%s://%s Server start error: undefine proto ", u.proto, u.addr))
	}
}

func (u *Server) Stop() {
	u.server.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("REQ: ", r.URL)
	splitArray := strings.Split(r.URL.String(), "/")
	if len(splitArray) == 0 {
		fmt.Fprintf(w, "need parameters")
		return
	}

	if strings.ToUpper(splitArray[1][0:1]) == "V" {
		if err := parseVersion(splitArray[1][1:]); err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	}
	switch {
	case strings.Contains(strings.ToLower(r.URL.String()), "createservice"):
		CreateService(w, r)
	case strings.Contains(strings.ToLower(r.URL.String()), "updateservice"):
		UpdateService(w, r)
	case strings.Contains(strings.ToLower(r.URL.String()), "removeservice"):
		RemoveService(w, r)
	case strings.Contains(strings.ToLower(r.URL.String()), "getserviceinfo"):
		GetServiceInfo(w, r)
	case strings.Contains(strings.ToLower(r.URL.String()), "gettaskinfo"):
		GetTaskInfo(w, r)
	case strings.Contains(strings.ToLower(r.URL.String()), "getnodeinfo"):
		GetNodeInfo(w, r)
	case strings.Contains(strings.ToLower(r.URL.String()), "getservicestate"):
		GetServiceState(w, r)
	default:
		otherCommand(w, r)
	}
}
func otherCommand(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Warning("ioutil.ReadAll faild")
		return
	}
	defer r.Body.Close()
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
	res, err := httpc.Post("http://unix/"+r.URL.String(),
		"application/x-www-form-urlencoded",
		strings.NewReader(string(body)))
	if err != nil {
		logrus.Warning("Post /var/run/docker.sock timeout")
		return
	}
	defer res.Body.Close()
	strbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Warning("ioutil.ReadAll faild")
		return
	}
	fmt.Fprintf(w, string(strbody))
}
