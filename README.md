# Yandex.Cloud Docker machine driver

## Installation
### Go
```shell
# install latest (git) version of docker-machine-driver-yandex in your $GOPATH/bin (depends on Golang and docker-machine)
$ go get -u github.com/yandex-cloud/docker-machine-driver-yandex
```

### Binary

WIP

## Usage
```bash
$ docker-machine create \
  --driver yandex \
  --yandex-token=FRegxx408323539fsdgHhsakawe0 \
  some-machine
```

## Options

|Option Name                                              |Description                                        |Default Value             |
|---------------------------------------------------------|---------------------------------------------------|--------------------------|
| ``--yandex-cloud-id`` or ``$YC_CLOUD_ID``               | Cloud ID                                          |                          |
| ``--yandex-cores`` or ``$YC_CORES``                     | Count of virtual CPUs                             | 2                        |
| ``--yandex-core-fraction`` or ``$YC_CORE_FRACTION``     | Core fraction                                     | 100                      |
| ``--yandex-disk-size`` or ``$YC_DISK_SIZE``             | Disk size in gigabytes                            | 20                       |
| ``--yandex-disk-type`` or ``$YC_DISK_TYPE``             | Disk type, e.g. 'network-hdd'                     | network-hdd              |
| ``--yandex-endpoint`` or ``$YC_ENDPOINT``               | Yandex.Cloud API Endpoint                         | api.cloud.yandex.net:443 |
| ``--yandex-folder-id`` or ``$YC_FOLDER_ID``             | Folder ID                                         |                          |
| ``--yandex-image-family`` or ``$YC_IMAGE_FAMILY``       | Image family name to lookup image ID for instance | ubuntu-1604-lts          |
| ``--yandex-image-folder-id`` or ``$YC_IMAGE_FOLDER_ID`` | Folder ID to the latest image by family name      | standard-images          |
| ``--yandex-image-id`` or ``$YC_IMAGE_ID``               | User-defined Image ID                             |                          |
| ``--yandex-labels`` or ``$YC_LABELS``                   | Instance labels in 'key=value' format             |                          |
| ``--yandex-memory`` or ``$YC_MEMORY``                   | Memory in gigabytes                               | 1                        |
| ``--yandex-nat`` or ``$YC_NAT``                         | Assign external (NAT) IP address                  | false                    |
| ``--yandex-platform-id`` or ``$YC_PLATFORM_ID``         | ID of the hardware platform configuration         | standard-v1              |
| ``--yandex-preemptible`` or ``$YC_PREEMPTIBLE``         | Yandex.Cloud Instance preemptibility flag         | false                    |
| ``--yandex-sa-key-file`` or ``$YC_SA_KEY_FILE``         | Yandex.Cloud Service Account key file             |                          |
| ``--yandex-ssh-port`` or ``$YC_SSH_PORT``               | SSH port                                          | 22                       |
| ``--yandex-ssh-user`` or ``$YC_SSH_USER``               | SSH username                                      | yc-user                  |
| ``--yandex-subnet-id`` or ``$YC_SUBNET_ID``             | Subnet ID                                         |                          |
| ``--yandex-token`` or ``$YC_TOKEN``                     | Yandex.Cloud OAuth token or IAM token             |                          |
| ``--yandex-use-internal-ip`` or ``$YC_USE_INTERNAL_IP`` | Use the internal Instance IP to communicate       | false                    |
| ``--yandex-userdata`` or ``$YC_USERDATA``               | Path to file with cloud-init user-data            |                          |
| ``--yandex-zone`` or ``$YC_ZONE``                       | Yandex.Cloud zone                                 | ru-central1-a            |
| ``--yandex-static-address`` or ``$YC_STATIC_ADDRESS``   | Set public static IPv4 address                                |                          |
---
