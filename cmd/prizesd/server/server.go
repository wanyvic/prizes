// +build linux freebsd openbsd darwin

package server // import "github.com/wanyvic/prizes/cmd/prizesd/server"

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/wanyvic/prizes/api"
	"github.com/wanyvic/prizes/cmd"
)

const DefaultProto = "unix"
const DefaultAddr = "./prizes.sock"
const DefaultHTTPHost = "localhost"
const DefaultHTTPPort = 9333 // Default HTTP Port
var DefaultTCPHost = fmt.Sprintf("%s:%d", DefaultHTTPHost, DefaultHTTPPort)

type ServerOpts struct {
	Proto string
	Addr  string
}
type Server struct {
	version string
	proto   string
	addr    string
	server  *http.Server
}

func (u *Server) setUnixDefaultOpts() {
	u.version = api.DefaultVersion
	u.proto = DefaultProto
	u.addr = DefaultAddr
}
func (u *Server) setTcpDefaultOpts() {
	u.version = api.DefaultVersion
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
		_ = server.Stop()
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
		_ = server.Stop()
	}()
	return &server, nil
}
func (u *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	u.server = &http.Server{
		Handler: mux,
	}
	if strings.Contains(u.proto, "unix") {
		os.Remove(u.addr)

		unixListener, err := net.Listen("unix", u.addr)
		if err != nil {
			panic(err)
		}
		u.server.Serve(unixListener)
	} else if strings.Contains(u.proto, "tcp") {
		unixListener, err := net.Listen("tcp", u.addr)
		if err != nil {
			panic(err)
		}
		u.server.Serve(unixListener)
	} else {
		return errors.New("undefine proto")
	}
	return nil
}

func (u *Server) Stop() error {

	u.server.Close()
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQ: ", r.URL)
	if strings.Contains(r.URL.String(), "CreateService") {
		CreateService(w, r)
	} else if strings.Contains(r.URL.String(), "UpdateService") {
		UpdateService(w, r)
	} else if strings.Contains(r.URL.String(), "RemoveService") {
		RemoveService(w, r)
	} else if strings.Contains(r.URL.String(), "GetServiceInfo") {
		GetServiceInfo(w, r)
	} else if strings.Contains(r.URL.String(), "GetTaskInfo") {
		GetTaskInfo(w, r)
	} else if strings.Contains(r.URL.String(), "GetNodeInfo") {
		GetNodeInfo(w, r)
	} else if strings.Contains(r.URL.String(), "GetServiceState") {
		GetServiceState(w, r)
	}
}
func CreateService(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll")
		return
	}
	var serviceSpec swarm.ServiceSpec
	if err := json.Unmarshal(body, &serviceSpec); err != nil {
		fmt.Fprintf(w, "bad parameters")
		return
	}
	response, err := cmd.CreateService(serviceSpec, types.ServiceCreateOptions{})
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(w, "json.Marshal error")
		return
	}
	fmt.Fprintf(w, string(jsonResponse))
}
func UpdateService(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQ: ", r.URL)
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]

	var serviceSpec swarm.ServiceSpec
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll")
		return
	}
	if err := json.Unmarshal(body, &serviceSpec); err != nil {
		fmt.Fprintf(w, "bad parameters")
		return
	}
	response, err := cmd.UpdateService(serviceID, serviceSpec, types.ServiceUpdateOptions{})
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	jsonResponse, err := json.Marshal(*response)
	if err != nil {
		fmt.Fprintf(w, "json.Marshal error")
		return
	}
	fmt.Fprintf(w, string(jsonResponse))
}
func RemoveService(w http.ResponseWriter, r *http.Request) {

	fmt.Println("REQ: ", r.URL)
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	err := cmd.RemoveService(serviceID)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "ServiceRemove Successed")
}
func GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQ: ", r.URL)
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	service, err := cmd.ServiceInfo(serviceID)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	json, err := json.Marshal(*service)
	if err != nil {
		fmt.Fprintf(w, "json.Marshal error")
		return
	}
	fmt.Fprintf(w, string(json))
}
func GetTaskInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQ: ", r.URL)
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	taskList, err := cmd.TasksInfo(serviceID)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	json, err := json.Marshal(*taskList)
	if err != nil {
		fmt.Fprintf(w, "json.Marshal error")
		return
	}
	fmt.Fprintf(w, string(json))
}
func GetNodeInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQ: ", r.URL)
	NodeID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	node, err := cmd.GetNodeInfo(NodeID)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	json, err := json.Marshal(*node)
	if err != nil {
		fmt.Fprintf(w, "json.Marshal error")
		return
	}
	fmt.Fprintf(w, string(json))
}
func GetServiceState(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQ: ", r.URL)
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	serviceStatistics, err := cmd.ServiceState(serviceID)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	json, err := json.Marshal(serviceStatistics)
	if err != nil {
		fmt.Fprintf(w, "json.Marshal error")
		return
	}
	fmt.Fprintf(w, string(json))
}
