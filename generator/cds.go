package generator

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/cluster"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"

	config "github.com/cookpad/itacho/api/v1/config"
)

// ConvertServiceDefinitionToCdsResources returns resources of Envoy xDS API.
func ConvertServiceDefinitionToCdsResources(def *config.ServiceDefinition, opts Opts) (*[]proto.Message, error) {
	deps := def.GetDependencies()
	rs := make([]proto.Message, len(deps))
	for i, dep := range deps {
		r, err := convertToCluster(dep, opts)
		if err != nil {
			return nil, err
		}
		rs[i] = r
	}
	return &rs, nil
}

func convertToCluster(dep *config.Dependency, opts Opts) (*api.Cluster, error) {
	var c *api.Cluster
	sds := dep.GetSds()
	if sds != nil && sds.Value {
		cl, err := convertToEdsCluster(dep, opts)
		if err != nil {
			return nil, err
		}
		c = cl
	} else {
		cl, err := convertToNonEdsCluster(dep)
		if err != nil {
			return nil, err
		}
		c = cl
	}

	br := dep.GetCircuitBreaker()
	if br != nil {
		c.CircuitBreakers = &cluster.CircuitBreakers{
			Thresholds: []*cluster.CircuitBreakers_Thresholds{
				&cluster.CircuitBreakers_Thresholds{
					Priority:           core.RoutingPriority_DEFAULT,
					MaxConnections:     &types.UInt32Value{Value: br.GetMaxConnections()},
					MaxPendingRequests: &types.UInt32Value{Value: br.GetMaxPendingRequests()},
					MaxRetries:         &types.UInt32Value{Value: br.GetMaxRetries()},
				},
			},
		}
	}

	detec := dep.GetOutlierDetection()
	if detec != nil {
		c.OutlierDetection = &cluster.OutlierDetection{
			Consecutive_5Xx: &types.UInt32Value{Value: detec.GetConsecutive_5Xx()},
		}
	}

	return c, nil
}

func convertToEdsCluster(dep *config.Dependency, opts Opts) (*api.Cluster, error) {
	delay := time.Duration(5 * time.Second)
	apiType := core.ApiConfigSource_REST
	if opts.LegacySds {
		apiType = core.ApiConfigSource_REST_LEGACY
	}

	c := api.Cluster{
		Name:              dep.GetClusterName(),
		Type:              api.Cluster_EDS,
		ConnectTimeout:    time.Duration(dep.GetConnectTimeoutMs()*1000*1000) * time.Nanosecond,
		LbPolicy:          api.Cluster_ROUND_ROBIN,
		ProtocolSelection: api.Cluster_USE_DOWNSTREAM_PROTOCOL,
		EdsClusterConfig: &api.Cluster_EdsClusterConfig{
			EdsConfig: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
					ApiConfigSource: &core.ApiConfigSource{
						ApiType:      apiType,
						ClusterNames: []string{opts.EdsCluster},
						RefreshDelay: &delay,
					},
				},
			},
		},
	}
	return &c, nil
}

func convertToNonEdsCluster(dep *config.Dependency) (*api.Cluster, error) {
	lbHost, port, err := extractHostAndPort(dep.GetLb())
	if err != nil {
		return nil, err
	}

	c := api.Cluster{
		Name:              dep.GetClusterName(),
		Type:              api.Cluster_STRICT_DNS,
		ConnectTimeout:    time.Duration(dep.GetConnectTimeoutMs()*1000*1000) * time.Nanosecond,
		LbPolicy:          api.Cluster_ROUND_ROBIN,
		ProtocolSelection: api.Cluster_USE_DOWNSTREAM_PROTOCOL,
		LoadAssignment: &api.ClusterLoadAssignment{
			ClusterName: dep.GetClusterName(),
			Endpoints: []endpoint.LocalityLbEndpoints{
				endpoint.LocalityLbEndpoints{
					LbEndpoints: []endpoint.LbEndpoint{*makeLbEndpoint(lbHost, port)},
				},
			},
		},
		// XXX: Just work-around for AWS EC2 environment, it could be configurable.
		DnsLookupFamily: api.Cluster_V4_ONLY,
	}
	tls := dep.GetTls()
	if tls != nil && tls.Value {
		c.TlsContext = &auth.UpstreamTlsContext{}

		sds := dep.GetSds()
		if sds == nil || !sds.Value {
			c.TlsContext.Sni = lbHost
		}
	}
	return &c, nil
}

func makeLbEndpoint(lbHost string, port uint32) *endpoint.LbEndpoint {
	return &endpoint.LbEndpoint{
		Endpoint: &endpoint.Endpoint{
			Address: &core.Address{
				Address: &core.Address_SocketAddress{
					SocketAddress: &core.SocketAddress{
						Address: lbHost,
						PortSpecifier: &core.SocketAddress_PortValue{
							PortValue: port,
						},
					},
				},
			},
		},
	}
}

func extractHostAndPort(lb string) (string, uint32, error) {
	s := strings.Split(lb, ":")
	if len(s) != 2 {
		return "", 0, fmt.Errorf("invalid lb value: %s", lb)
	}
	port, err := strconv.ParseUint(s[1], 10, 32)
	if err != nil {
		return "", 0, fmt.Errorf("invalid lb value: %s", lb)
	}
	return s[0], uint32(port), nil
}
