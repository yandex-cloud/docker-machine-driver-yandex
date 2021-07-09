module github.com/GennadySpb/docker-machine-driver-yandex

go 1.16

require (
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/docker/machine v0.16.2
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/stretchr/testify v1.5.1
	github.com/yandex-cloud/docker-machine-driver-yandex v0.1.27-rc
	github.com/yandex-cloud/go-genproto v0.0.0-20201102102956-0c505728b6f0
	github.com/yandex-cloud/go-sdk v0.0.0-20201109103511-a86298d3fea5
	google.golang.org/grpc v1.29.1
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200309214505-aa6a9891b09c

//replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.1
