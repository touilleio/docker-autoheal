module github.com/touilleio/docker-autoheal

go 1.15

require (
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v20.10.1+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.34.0 // indirect
	github.com/sqooba/go-common/logging v0.0.0-20201231000000
	github.com/sqooba/go-common/version v0.0.0-20201231000000
)

replace (
	//github.com/sqooba/go-common/healthchecks => ../sqooba-go-common/healthchecks
	github.com/sqooba/go-common/logging => ../sqooba-go-common/logging
	github.com/sqooba/go-common/version => ../sqooba-go-common/version
)
