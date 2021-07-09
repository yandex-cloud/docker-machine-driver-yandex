module github.com/GennadySpb/docker-machine-driver-yandex

go 1.16

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/machine v0.16.2
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/stretchr/testify v1.5.1
	github.com/yandex-cloud/go-genproto v0.0.0-20201102102956-0c505728b6f0
	github.com/yandex-cloud/go-sdk v0.0.0-20201109103511-a86298d3fea5
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200309214505-aa6a9891b09c

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.1
