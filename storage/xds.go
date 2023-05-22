package storage

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/cookpad/itacho/xds"
)

// XdsResponseStorage is a gateway client for object storage.
type XdsResponseStorage interface {
	Fetch(t xds.ResponseType, v xds.APIVersion, nodeCluster string) (*int, *[]byte, error)
}

// NewXdsResponseStorage returns configured XdsResponseStorage.
func NewXdsResponseStorage(endpoint url.URL) XdsResponseStorage {
	return &objectStorageGateway{endpoint}
}

type objectStorageGateway struct {
	endpoint url.URL
}

func versionPrefix(v xds.APIVersion) string {
	switch v {
	case xds.V2:
		return "/v2"
	case xds.V3:
		return "/v3"
	default:
		panic("unsupported API version")
	}
}

func cdsPathForNode(v xds.APIVersion, nodeCluster string) string {
	return filepath.Join(versionPrefix(v), "discovery/clusters", nodeCluster)
}

func rdsPathForNode(v xds.APIVersion, nodeCluster string) string {
	return filepath.Join(versionPrefix(v), "discovery/routes", nodeCluster)
}

// Fetch xDS response JSON from object storage.
func (s *objectStorageGateway) Fetch(t xds.ResponseType, v xds.APIVersion, nodeCluster string) (*int, *[]byte, error) {
	u := s.endpoint

	switch t {
	case xds.CDS:
		{
			u.Path = cdsPathForNode(v, nodeCluster)
		}
	case xds.RDS:
		{
			u.Path = rdsPathForNode(v, nodeCluster)
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
