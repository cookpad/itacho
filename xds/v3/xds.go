package xds

import (
	"bytes"
	"fmt"

	"github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
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
func ExtractNodeCluster(node *envoy_config_core_v3.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Cluster
}

// UnmarshalDiscoveryRequest build Envoy's DiscoveryRequest proto message from JSON string
func UnmarshalDiscoveryRequest(typeURL string, body *[]byte) (*envoy_service_discovery_v3.DiscoveryRequest, error) {
	req := &envoy_service_discovery_v3.DiscoveryRequest{}
	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err := unmarshaler.Unmarshal(bytes.NewReader(*body), req); err != nil {
		return nil, fmt.Errorf("Failed parse JSON body: %s", err)
	}
	req.TypeUrl = typeURL
	return req, nil
}
