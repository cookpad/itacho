package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/cookpad/itacho/storage"
	"github.com/cookpad/itacho/xds"
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
	switch p {
	case "/v2/discovery:clusters":
		typeURL = xds.ClusterType
		t = xds.CDS
	case "/v2/discovery:routes":
		typeURL = xds.RouteType
		t = xds.RDS
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
	xdsReq, err := xds.UnmarshalDiscoveryRequest(typeURL, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nodeCluster := xds.ExtractNodeCluster(xdsReq.GetNode())
	log.Infof("nodeCluster=%s", nodeCluster)

	code, jsonStr, err := h.storage.Fetch(t, nodeCluster)
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
