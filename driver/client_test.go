package driver

import (
	"testing"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/require"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
)

func TestNewYandexCloudClient(t *testing.T) {
	type args struct {
		d *Driver
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "use token",
			args: args{
				d: &Driver{
					Token:                 "some-test-token",
					ServiceAccountKeyFile: "",
				},
			},
			wantErr: false,
		},
		{
			name: "service account key file doest not exist",
			args: args{
				d: &Driver{
					Token:                 "",
					ServiceAccountKeyFile: "some-not-exist-sa-key-file",
				},
			},
			wantErr: true,
		},
		{
			name: "both auth methods provided",
			args: args{
				d: &Driver{
					Token:                 "some-test-token",
					ServiceAccountKeyFile: "testdata/fake_service_account_key.json",
				},
			},
			wantErr: true,
		},
		{
			name: "service account key file",
			args: args{
				d: &Driver{
					Token:                 "",
					ServiceAccountKeyFile: "testdata/fake_service_account_key.json",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewYCClient(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewYCClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestYandexCloudClient_getInstanceIPAddress(t *testing.T) {
	type args struct {
		d        *Driver
		instance *compute.Instance
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string
		wantErr     bool
	}{
		{
			name: "instance with both addresses, want internal address",
			args: args{
				d: &Driver{
					UseIPv6:       false,
					UseInternalIP: true,
				},
				instance: &compute.Instance{
					NetworkInterfaces: []*compute.NetworkInterface{
						{
							Index: "1",
							PrimaryV4Address: &compute.PrimaryAddress{
								Address: "192.168.19.86",
								OneToOneNat: &compute.OneToOneNat{
									Address:   "92.68.12.34",
									IpVersion: compute.IpVersion_IPV4,
								},
							},
							SubnetId:   "some-subnet-id",
							MacAddress: "aa-bb-cc-dd-ee-ff",
						},
					},
				},
			},
			wantAddress: "192.168.19.86",
			wantErr:     false,
		},
		{
			name: "instance with both addresses, want external address",
			args: args{
				d: &Driver{
					UseIPv6:       false,
					UseInternalIP: false,
				},
				instance: &compute.Instance{
					NetworkInterfaces: []*compute.NetworkInterface{
						{
							Index: "1",
							PrimaryV4Address: &compute.PrimaryAddress{
								Address: "192.168.19.86",
								OneToOneNat: &compute.OneToOneNat{
									Address:   "92.68.12.34",
									IpVersion: compute.IpVersion_IPV4,
								},
							},
							SubnetId:   "some-subnet-id",
							MacAddress: "aa-bb-cc-dd-ee-ff",
						},
					},
				},
			},
			wantAddress: "92.68.12.34",
			wantErr:     false,
		},
		{
			name: "instance with internal address, want external address",
			args: args{
				d: &Driver{
					UseIPv6:       false,
					UseInternalIP: false,
				},
				instance: &compute.Instance{
					NetworkInterfaces: []*compute.NetworkInterface{
						{
							Index: "1",
							PrimaryV4Address: &compute.PrimaryAddress{
								Address: "192.168.19.86",
							},
							SubnetId:   "some-subnet-id",
							MacAddress: "aa-bb-cc-dd-ee-ff",
						},
					},
				},
			},
			wantAddress: "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &YCClient{}
			gotAddress, err := c.getInstanceIPAddress(tt.args.d, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("YCClient.getInstanceIPAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("YCClient.getInstanceIPAddress() = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}

func TestYandexCloudClient_instanceAddresses(t *testing.T) {
	type args struct {
		instance *compute.Instance
	}
	tests := []struct {
		name        string
		args        args
		wantIPV4Int string
		wantIPV4Ext string
		wantIPV6    string
		wantErr     bool
	}{
		{
			name: "no nics defined",
			args: args{
				instance: &compute.Instance{
					NetworkInterfaces: []*compute.NetworkInterface{},
				},
			},
			wantIPV4Int: "",
			wantIPV4Ext: "",
			wantIPV6:    "",
			wantErr:     true,
		},
		{
			name: "one nic with internal address",
			args: args{
				instance: &compute.Instance{
					NetworkInterfaces: []*compute.NetworkInterface{
						{
							Index: "1",
							PrimaryV4Address: &compute.PrimaryAddress{
								Address: "192.168.19.16",
							},
							SubnetId:   "some-subnet-id",
							MacAddress: "aa-bb-cc-dd-ee-ff",
						},
					},
				},
			},
			wantIPV4Int: "192.168.19.16",
			wantIPV4Ext: "",
			wantIPV6:    "",
			wantErr:     false,
		},
		{
			name: "one nic with internal and external address",
			args: args{
				instance: &compute.Instance{
					NetworkInterfaces: []*compute.NetworkInterface{
						{
							Index: "1",
							PrimaryV4Address: &compute.PrimaryAddress{
								Address: "192.168.19.86",
								OneToOneNat: &compute.OneToOneNat{
									Address:   "92.68.12.34",
									IpVersion: compute.IpVersion_IPV4,
								},
							},
							SubnetId:   "some-subnet-id",
							MacAddress: "aa-bb-cc-dd-ee-ff",
						},
					},
				},
			},
			wantIPV4Int: "192.168.19.86",
			wantIPV4Ext: "92.68.12.34",
			wantIPV6:    "",
			wantErr:     false,
		},
		{
			name: "one nic with ipv6 address",
			args: args{
				instance: &compute.Instance{
					NetworkInterfaces: []*compute.NetworkInterface{
						{
							Index: "1",
							PrimaryV4Address: &compute.PrimaryAddress{
								Address: "192.168.19.86",
							},
							PrimaryV6Address: &compute.PrimaryAddress{
								Address: "2001:db8::370:7348",
							},
							SubnetId:   "some-subnet-id",
							MacAddress: "aa-bb-cc-dd-ee-ff",
						},
					},
				},
			},
			wantIPV4Int: "192.168.19.86",
			wantIPV4Ext: "",
			wantIPV6:    "2001:db8::370:7348",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &YCClient{}
			gotIPV4Int, gotIPV4Ext, gotIPV6, err := c.instanceAddresses(tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("YCClient.instanceAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIPV4Int != tt.wantIPV4Int {
				t.Errorf("YCClient.instanceAddresses() gotIPV4Int = %v, want %v", gotIPV4Int, tt.wantIPV4Int)
			}
			if gotIPV4Ext != tt.wantIPV4Ext {
				t.Errorf("YCClient.instanceAddresses() gotIPV4Ext = %v, want %v", gotIPV4Ext, tt.wantIPV4Ext)
			}
			if gotIPV6 != tt.wantIPV6 {
				t.Errorf("YCClient.instanceAddresses() gotIPV6 = %v, want %v", gotIPV6, tt.wantIPV6)
			}
		})
	}
}

func Test_prepareInstanceCreateRequest(t *testing.T) {
	type args struct {
		d       *Driver
		imageID string
	}
	var tests = []struct {
		name string
		args args
		want *compute.CreateInstanceRequest
	}{
		{
			name: "a typical instance without nat",
			args: args{
				d: &Driver{
					BaseDriver: &drivers.BaseDriver{
						MachineName: "foobar-name",
					},
					Cores:         2,
					CoreFraction:  100,
					DiskSize:      20,
					DiskType:      "network-hdd",
					FolderID:      "some-folder-id",
					Memory:        2,
					Nat:           false,
					PlatformID:    "standard-v2",
					Preemptible:   false,
					SubnetID:      "foobar-subnet",
					UseIPv6:       false,
					UseInternalIP: false,
					Zone:          "ru-central1-c",
				},
				imageID: "foobar-image-id",
			},
			want: &compute.CreateInstanceRequest{
				FolderId:    "some-folder-id",
				Name:        "foobar-name",
				Description: "",
				Labels:      map[string]string{},
				ZoneId:      "ru-central1-c",
				PlatformId:  "standard-v2",
				ResourcesSpec: &compute.ResourcesSpec{
					Memory:       toBytes(2),
					Cores:        2,
					CoreFraction: 100,
					Gpus:         0,
				},
				BootDiskSpec: &compute.AttachedDiskSpec{
					AutoDelete: true,
					Disk: &compute.AttachedDiskSpec_DiskSpec_{
						DiskSpec: &compute.AttachedDiskSpec_DiskSpec{
							TypeId: "network-hdd",
							Size:   toBytes(20),
							Source: &compute.AttachedDiskSpec_DiskSpec_ImageId{
								ImageId: "foobar-image-id",
							},
						},
					},
				},
				SecondaryDiskSpecs: nil,
				NetworkInterfaceSpecs: []*compute.NetworkInterfaceSpec{
					{
						SubnetId:             "foobar-subnet",
						PrimaryV4AddressSpec: &compute.PrimaryAddressSpec{},
						PrimaryV6AddressSpec: nil,
						SecurityGroupIds:     nil,
					},
				},
				Hostname:         "",
				ServiceAccountId: "",
				SchedulingPolicy: &compute.SchedulingPolicy{},
			},
		},
		{
			name: "instance with nat and ssd disk and user-data",
			args: args{
				d: &Driver{
					BaseDriver: &drivers.BaseDriver{
						MachineName: "foobar-name",
					},
					Cores:         2,
					CoreFraction:  50,
					DiskSize:      20,
					DiskType:      "network-ssd",
					FolderID:      "some-folder-id",
					Memory:        2,
					Nat:           true,
					PlatformID:    "standard-v2",
					Preemptible:   false,
					SubnetID:      "foobar-subnet",
					UseInternalIP: true,
					Zone:          "ru-central1-c",
				},
				imageID: "foobar-image-id",
			},
			want: &compute.CreateInstanceRequest{
				FolderId:    "some-folder-id",
				Name:        "foobar-name",
				Description: "",
				Labels:      map[string]string{},
				ZoneId:      "ru-central1-c",
				PlatformId:  "standard-v2",
				ResourcesSpec: &compute.ResourcesSpec{
					Memory:       toBytes(2),
					Cores:        2,
					CoreFraction: 50,
					Gpus:         0,
				},
				BootDiskSpec: &compute.AttachedDiskSpec{
					AutoDelete: true,
					Disk: &compute.AttachedDiskSpec_DiskSpec_{
						DiskSpec: &compute.AttachedDiskSpec_DiskSpec{
							TypeId: "network-ssd",
							Size:   toBytes(20),
							Source: &compute.AttachedDiskSpec_DiskSpec_ImageId{
								ImageId: "foobar-image-id",
							},
						},
					},
				},
				SecondaryDiskSpecs: nil,
				NetworkInterfaceSpecs: []*compute.NetworkInterfaceSpec{
					{
						SubnetId: "foobar-subnet",
						PrimaryV4AddressSpec: &compute.PrimaryAddressSpec{
							OneToOneNatSpec: &compute.OneToOneNatSpec{
								IpVersion: compute.IpVersion_IPV4,
							},
						},
						PrimaryV6AddressSpec: nil,
						SecurityGroupIds:     nil,
					},
				},
				Hostname:         "",
				ServiceAccountId: "",
				SchedulingPolicy: &compute.SchedulingPolicy{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prepareInstanceCreateRequest(tt.args.d, tt.args.imageID)
			require.Equal(t, tt.want, got)
		})
	}
}
