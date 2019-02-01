package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	gogojsonpb "github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/jsonpb"
	jsonnet "github.com/google/go-jsonnet"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	config "github.com/cookpad/itacho/api/v1/config"
	"github.com/cookpad/itacho/xds"
)

// Opts is a set of options for generator
type Opts struct {
	SourcePath  string
	OutputDir   string
	Type        xds.ResponseType
	NodeCluster string
	Version     string
	LegacySds   bool
	EdsCluster  string
}

// Generate xDS response flagment.
func Generate(opts Opts) error {
	path := opts.SourcePath
	nodeCluster := extractNodeCluster(path)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf(`Failed to open file "%s": %s`, path, err)
	}

	vm := jsonnet.MakeVM()
	jsonStr, err := vm.EvaluateSnippet(path, string(content))
	if err != nil {
		return fmt.Errorf(`Failed to evaluate Jsonnet file "%s": %s`, path, err)
	}

	def := &config.ServiceDefinition{}
	// XXX: Using gogoproto's jsonpb causes an error (can not handle oneof field)
	if err := jsonpb.UnmarshalString(jsonStr, def); err != nil {
		return fmt.Errorf("Failed to load json: %s", err)
	}
	if err := def.Validate(); err != nil {
		return fmt.Errorf("Failed to validate ServiceDefinition: %s", err)
	}

	res, err := convertToXdsRespose(opts.Type, opts.Version, def, opts)
	if err != nil {
		return fmt.Errorf("Failed to convert to xDS response: %s", err)
	}

	buf := &bytes.Buffer{}
	if err := (&gogojsonpb.Marshaler{OrigName: true}).Marshal(buf, res); err != nil {
		return fmt.Errorf("Failed to marshal protos: %s", err)
	}

	outPath := outputPath(opts.Type, opts.OutputDir, nodeCluster)
	if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create output directories: %s", err)
	}
	if err := ioutil.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("Failed to write response flagment: %s", err)
	}
	log.Infof("Generated xDS response: %s", outPath)

	return nil
}

// GenerateYaml writes evaluated Jsonnet content into `output` as YAML format.
func GenerateYaml(source string, output string) error {
	content, err := ioutil.ReadFile(source)
	if err != nil {
		return fmt.Errorf(`Failed to open file "%s": %s`, source, err)
	}

	vm := jsonnet.MakeVM()
	jsonStr, err := vm.EvaluateSnippet(source, string(content))
	if err != nil {
		return fmt.Errorf(`Failed to evaluate Jsonnet file "%s": %s`, source, err)
	}

	m := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(jsonStr), &m); err != nil {
		return fmt.Errorf("Failed to Unmarshal json string: %s", err)
	}
	yamlStr, err := yaml.Marshal(&m)
	if err != nil {
		return fmt.Errorf("Failed to marshal hash map into YAML: %s", err)
	}

	if err := os.MkdirAll(filepath.Dir(output), os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create output directories: %s", err)
	}
	if err := ioutil.WriteFile(output, []byte(yamlStr), 0644); err != nil {
		return fmt.Errorf("Failed to write YAML content: %s", err)
	}

	log.Infof("Generated YAML file: %s", output)
	return nil
}

func convertToXdsRespose(t xds.ResponseType, version string, def *config.ServiceDefinition, opts Opts) (*v2.DiscoveryResponse, error) {
	var res *v2.DiscoveryResponse

	switch t {
	case xds.CDS:
		{
			cs, err := ConvertServiceDefinitionToCdsResources(def, opts)
			if err != nil {
				return nil, err
			}
			out, err := createResponse(cs, xds.ClusterType, version)
			if err != nil {
				return nil, err
			}
			res = out
		}
	case xds.RDS:
		{
			rs, err := ConvertServiceDefinitionToRdsResources(def)
			if err != nil {
				return nil, err
			}
			out, err := createResponse(rs, xds.RouteType, version)
			if err != nil {
				return nil, err
			}
			res = out
		}
	}

	return res, nil
}

func createResponse(protos *[]proto.Message, typeURL string, version string) (*v2.DiscoveryResponse, error) {
	resources := make([]types.Any, len(*protos))
	for i := 0; i < len(*protos); i++ {
		data, err := proto.Marshal((*protos)[i])
		if err != nil {
			return nil, err
		}
		resources[i] = types.Any{
			TypeUrl: typeURL,
			Value:   data,
		}
	}
	return &v2.DiscoveryResponse{
		VersionInfo: version,
		Resources:   resources,
		TypeUrl:     typeURL,
	}, nil
}

func outputPath(t xds.ResponseType, dir string, nodeCluster string) string {
	var out string
	switch t {
	case xds.CDS:
		{
			out = filepath.Join(dir, xds.CdsPathForNode(nodeCluster))
		}
	case xds.RDS:
		{
			out = filepath.Join(dir, xds.RdsPathForNode(nodeCluster))
		}
	}
	return out
}

func extractNodeCluster(path string) string {
	basePath := filepath.Base(path)
	return strings.TrimSuffix(basePath, filepath.Ext(basePath))
}
