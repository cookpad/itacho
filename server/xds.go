package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/cookpad/itacho/storage"
	"github.com/cookpad/itacho/xds"
	v2xds "github.com/cookpad/itacho/xds/v2"
	v3xds "github.com/cookpad/itacho/xds/v3"
)

// XdsHandler for http requests.
type XdsHandler interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// NewXdsHandler returns new Handler instance.
func NewXdsHandler(str storage.XdsResponseStorage) XdsHandler {
	return &xdsHandler{str}
}

type xdsHandler struct {
	storage storage.XdsResponseStorage
}

func (h *xdsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p := path.Clean(req.URL.Path)

	var typeURL string
	var t xds.ResponseType
	var v xds.APIVersion
	switch p {
	case "/v2/discovery:clusters":
		typeURL = v2xds.ClusterType
		t = xds.CDS
		v = xds.V2
	case "/v2/discovery:routes":
		typeURL = v2xds.RouteType
		t = xds.RDS
		v = xds.V2
	case "/v3/discovery:clusters":
		typeURL = v3xds.ClusterType
		t = xds.CDS
		v = xds.V3
	case "/v3/discovery:routes":
		typeURL = v3xds.RouteType
		t = xds.RDS
		v = xds.V3
	default:
		http.Error(w, "no endpoint", http.StatusNotFound)
		return
	}

	if req.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	var nodeCluster string
	switch v {
	case xds.V2:
		xdsReq, err := v2xds.UnmarshalDiscoveryRequest(typeURL, &body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nodeCluster = v2xds.ExtractNodeCluster(xdsReq.GetNode())
	case xds.V3:
		xdsReq, err := v3xds.UnmarshalDiscoveryRequest(typeURL, &body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nodeCluster = v3xds.ExtractNodeCluster(xdsReq.GetNode())
	}

	log.Infof("nodeCluster=%s", nodeCluster)
	code, jsonStr, err := h.storage.Fetch(t, v, nodeCluster)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get content from object storage: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(*code)
	if _, err = w.Write(*jsonStr); err != nil {
		log.Errorf("gateway error: %v", err)
	}
	return
}
