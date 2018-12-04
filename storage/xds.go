package storage

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cookpad/itacho/xds"
)

// XdsResponseStorage is a gateway client for object storage.
type XdsResponseStorage interface {
	Fetch(t xds.ResponseType, nodeCluster string) (*int, *[]byte, error)
}

// NewXdsResponseStorage returns configured XdsResponseStorage.
func NewXdsResponseStorage(endpoint url.URL) XdsResponseStorage {
	return &objectStorageGateway{endpoint}
}

type objectStorageGateway struct {
	endpoint url.URL
}

// Fetch xDS response JSON from object storage.
func (s *objectStorageGateway) Fetch(t xds.ResponseType, nodeCluster string) (*int, *[]byte, error) {
	u := s.endpoint

	switch t {
	case xds.CDS:
		{
			u.Path = xds.CdsPathForNode(nodeCluster)
		}
	case xds.RDS:
		{
			u.Path = xds.RdsPathForNode(nodeCluster)
		}
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return &resp.StatusCode, &body, nil
}
