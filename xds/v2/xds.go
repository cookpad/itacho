package xds

import (
	"fmt"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v2"
	"github.com/golang/protobuf/jsonpb"
)

const (
	// EndpointType for xDS resource type
	EndpointType = resource.EndpointType
	// ClusterType for xDS resource type
	ClusterType = resource.ClusterType
	// RouteType for xDS resource type
	RouteType = resource.RouteType
)

// ExtractNodeCluster returns cluster value from Node proto message
func ExtractNodeCluster(node *envoy_api_v2_core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Cluster
}

// UnmarshalDiscoveryRequest build Envoy's DiscoveryRequest proto message from JSON string
func UnmarshalDiscoveryRequest(typeURL string, body *[]byte) (*envoy_api_v2.DiscoveryRequest, error) {
	req := &envoy_api_v2.DiscoveryRequest{}
	if err := jsonpb.UnmarshalString(string(*body), req); err != nil {
		return nil, fmt.Errorf("Failed parse JSON body: %s", err)
	}
	req.TypeUrl = typeURL
	return req, nil
}
