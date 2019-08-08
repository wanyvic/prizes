package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api"
	"github.com/wanyvic/prizes/api/types/service"
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
func ServiceCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Warning("ioutil.ReadAll faild")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, parseError(errors.New("parameters invalid")))
		return
	}
	logrus.Debug("body: ", string(body))
	defer r.Body.Close()
	serviceCreate := service.ServiceCreate{}
	if err := json.Unmarshal(body, &serviceCreate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, parseError(errors.New("parameters invalid")))
		return
	}
	response, err := cmd.ServiceCreate(&serviceCreate)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(response)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	logrus.Info(fmt.Sprintf("http response ID: %s ,Warning: %s", response.ID, response.Warnings))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, strResult)
}
func ServiceUpdate(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, parseError(errors.New("parameters invalid")))
		return
	}
	logrus.Debug("body: ", string(body))
	defer r.Body.Close()
	serviceUpdate := service.ServiceUpdate{}
	if err := json.Unmarshal(body, &serviceUpdate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, parseError(errors.New("parameters invalid")))
		return
	}
	if serviceID != serviceUpdate.ServiceID {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintf(w, parseError(fmt.Errorf("serviceID mismatch "+serviceID+" "+serviceUpdate.ServiceID)))
		return
	}
	response, err := cmd.ServiceUpdate(&serviceUpdate)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(response)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	logrus.Info(fmt.Sprintf("http response ID: %s ,Warning: %s", serviceID, response.Warnings))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strResult)
}
func ServiceStatement(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	statement, err := cmd.ServiceStatement(serviceID, time.Now().UTC())
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(*statement)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strResult)
}
func ServiceRefund(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	refunInfo, err := cmd.ServiceRefund(serviceID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(*refunInfo)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, strResult)
}
func GetService(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	serviceInfo, err := cmd.ServiceInfo(serviceID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(*serviceInfo)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strResult)
}
func GetServicesFromPubkey(w http.ResponseWriter, r *http.Request) {
	pubkey := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	fmt.Println(pubkey)
	serviceInfoList, err := cmd.GetServicesFromPubkey(pubkey)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(*serviceInfoList)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strResult)
}
func GetNode(w http.ResponseWriter, r *http.Request) {
	NodeID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	node, err := cmd.GetNodeInfo(NodeID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(*node)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strResult)
}
func GetNodeList(w http.ResponseWriter, r *http.Request) {
	nodeListStatistics, err := cmd.GetNodeList()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(*nodeListStatistics)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strResult)
}
func GetServiceState(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.String()[strings.LastIndex(r.URL.String(), "/")+1:]
	serviceStatistics, err := cmd.ServiceState(serviceID, time.Now().UTC())
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, parseError(err))
		return
	}
	strResult, successd := parseResult(serviceStatistics)
	if !successd {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, strResult)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strResult)
}
