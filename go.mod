module github.com/yandex-cloud/docker-machine-driver-yandex

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/machine v0.16.2
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/yandex-cloud/go-genproto v0.0.0-20200316083924-bd7a668b4f7b
	github.com/yandex-cloud/go-sdk v0.0.0-20200306133551-1e4c0ef0a537
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4 // indirect
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200309214505-aa6a9891b09c
