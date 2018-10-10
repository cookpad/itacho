itacho
======

itacho _板長_ to manage and operate envoy based service mesh.

## Configuration

Environment variables

- `OBJECT_STORAGE_ENDPOINT_URL`: an endpoint for object storage
- `BIND_PORT`: [optional] a port number to bind and listen

For further detail, see `itacho --help` and `itacho ${sub_cmd} --help`.

## Design notes
### Object storage path convention

- Cluster: `GET /v2/discovery/clusters/${node_cluster}`
- Route: `GET /v2/discovery/routes/${node_cluster}`

## Development
### build proto files

```
go get github.com/gogo/protobuf/proto
go get github.com/gogo/protobuf/gogoproto
go get github.com/gogo/protobuf/protoc-gen-gofast
go get github.com/lyft/protoc-gen-validate
go get github.com/goware/modvendor
```

```
make
make integration_test
```
