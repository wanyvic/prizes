package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api"
	"github.com/wanyvic/prizes/cmd"
)

func parseVersion(strVersion string) error {
	logrus.Info("ParseVersion")
	version, err := strconv.ParseFloat(strVersion, 64)
	if err != nil {
		return nil
	}
	apiVersion, _ := strconv.ParseFloat(api.MinAPIVersion, 64)
	if version < apiVersion {
		return fmt.Errorf("version %0.2f is old than api version %0.2f", version, apiVersion)
	}
	return nil
}
func CreateService(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Warning("ioutil.ReadAll faild")
		return
	}
	defer r.Body.Close()
	serviceSpec := swarm.ServiceSpec{}
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
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	serviceSpec := swarm.ServiceSpec{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Warning("ioutil.ReadAll faild")
		return
	}
	defer r.Body.Close()
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
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	err := cmd.RemoveService(serviceID)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "ServiceRemove Successed")
}
func GetServiceInfo(w http.ResponseWriter, r *http.Request) {
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
