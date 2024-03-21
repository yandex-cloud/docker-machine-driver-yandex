package driver

import (
	"context"
	"errors"
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/docker/machine/libmachine/log"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/pkg/requestid"
	"github.com/yandex-cloud/go-sdk/pkg/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const MaxRetries = 3
const StandardImagesFolderID = "standard-images"
const UserAgent = "docker-machine-driver-yandex"

type YCClient struct {
	sdk *ycsdk.SDK
}

func (c *YCClient) createInstance(d *Driver) error {
	ctx := context.Background()

	imageID := d.ImageID
	if imageID == "" {
		var err error
		imageID, err = c.getImageIDFromFolder(d.ImageFamily, d.ImageFolderID)
		if err != nil {
			return err
		}
	}

	log.Infof("Use image with ID %q from folder ID %q", imageID, d.ImageFolderID)

	request := prepareInstanceCreateRequest(d, imageID)

	op, err := c.sdk.WrapOperation(c.sdk.Compute().Instance().Create(ctx, request))
	if err != nil {
		return fmt.Errorf("Error while requesting API to create instance: %s", err)
	}

	protoMetadata, err := op.Metadata()
	if err != nil {
		return fmt.Errorf("Error while get instance create operation metadata: %s", err)
	}

	md, ok := protoMetadata.(*compute.CreateInstanceMetadata)
	if !ok {
		return fmt.Errorf("could not get Instance ID from create operation metadata")
	}

	d.InstanceID = md.InstanceId

	log.Infof("Waiting for Instance with ID %q", d.InstanceID)
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("Error while waiting operation to create instance: %s", err)
	}

	resp, err := op.Response()
	if err != nil {
		return fmt.Errorf("Instance creation failed: %s", err)
	}

	instance, ok := resp.(*compute.Instance)
	if !ok {
		return fmt.Errorf("Create response doesn't contain Instance")
	}

	d.IPAddress, err = c.getInstanceIPAddress(d, instance)

	return err
}

func prepareInstanceCreateRequest(d *Driver, imageID string) *compute.CreateInstanceRequest {
	// TODO support static address assignment
	// TODO additional disks

	request := &compute.CreateInstanceRequest{
		FolderId:   d.FolderID,
		Name:       d.MachineName,
		ZoneId:     d.Zone,
		PlatformId: d.PlatformID,
		ResourcesSpec: &compute.ResourcesSpec{
			Cores:        int64(d.Cores),
			Gpus:         int64(d.Gpus),
			CoreFraction: int64(d.CoreFraction),
			Memory:       toBytes(d.Memory),
		},
		BootDiskSpec: &compute.AttachedDiskSpec{
			AutoDelete: true,
			Disk: &compute.AttachedDiskSpec_DiskSpec_{
				DiskSpec: &compute.AttachedDiskSpec_DiskSpec{
					TypeId: d.DiskType,
					Size:   toBytes(d.DiskSize),
					Source: &compute.AttachedDiskSpec_DiskSpec_ImageId{
						ImageId: imageID,
					},
				},
			},
		},
		Labels: d.ParsedLabels(),
		NetworkInterfaceSpecs: []*compute.NetworkInterfaceSpec{
			{
				SubnetId:             d.SubnetID,
				PrimaryV4AddressSpec: &compute.PrimaryAddressSpec{},
				SecurityGroupIds:     d.SecurityGroups,
			},
		},
		SchedulingPolicy: &compute.SchedulingPolicy{
			Preemptible: d.Preemptible,
		},
		ServiceAccountId: d.ServiceAccountID,
		Metadata:         d.Metadata,
	}

	if d.Nat {
		if d.StaticAddress == "" {
			request.NetworkInterfaceSpecs[0].PrimaryV4AddressSpec.OneToOneNatSpec = &compute.OneToOneNatSpec{
				IpVersion: compute.IpVersion_IPV4,
			}
		} else {
			request.NetworkInterfaceSpecs[0].PrimaryV4AddressSpec.OneToOneNatSpec = &compute.OneToOneNatSpec{
				Address:   d.StaticAddress,
				IpVersion: compute.IpVersion_IPV4,
			}
		}
	}

	if len(d.Filesystems) > 0 {
		var fsSpecs []*compute.AttachedFilesystemSpec
		fs, err := d.ParseFilesystems()
		if err != nil {
			//If error set empty Filesystems
			log.Infof("Error in filesystem format %q", err)
			request.FilesystemSpecs = fsSpecs
		} else {
			for deviceName, fileSystem := range fs {
				fsSpecs = append(fsSpecs, &compute.AttachedFilesystemSpec{
					DeviceName:   deviceName,
					FilesystemId: fileSystem["filesystemId"],
					Mode:         compute.AttachedFilesystemSpec_READ_WRITE,
				})
			}
			request.FilesystemSpecs = fsSpecs
		}
	}

	return request
}

func NewYCClient(d *Driver) (*YCClient, error) {
	credentials, err := d.Credentials()
	if err != nil {
		return nil, err
	}

	config := ycsdk.Config{
		Credentials: credentials,
	}

	if d.Endpoint != "" {
		config.Endpoint = d.Endpoint
	}

	headerMD := metadata.Pairs("user-agent", UserAgent)

	requestIDInterceptor := requestid.Interceptor()

	retryInterceptor := retry.Interceptor(
		retry.WithMax(MaxRetries),
		retry.WithCodes(codes.Unavailable),
		retry.WithAttemptHeader(true),
		retry.WithBackoff(retry.DefaultBackoff()),
	)

	// Make sure retry interceptor is above id interceptor.
	// Now we will have new request id for every retry attempt.
	interceptorChain := grpc_middleware.ChainUnaryClient(retryInterceptor, requestIDInterceptor)

	ctxWithClientTraceID := requestid.ContextWithClientTraceID(context.Background(), uuid.New().String())
	sdk, err := ycsdk.Build(ctxWithClientTraceID, config,
		grpc.WithUserAgent(UserAgent),
		grpc.WithDefaultCallOptions(grpc.Header(&headerMD)),
		grpc.WithUnaryInterceptor(interceptorChain))

	if err != nil {
		return nil, err
	}

	return &YCClient{
		sdk: sdk,
	}, nil
}

func (c *YCClient) getImageIDFromFolder(familyName, lookupFolderID string) (string, error) {
	image, err := c.sdk.Compute().Image().GetLatestByFamily(context.Background(), &compute.GetImageLatestByFamilyRequest{
		FolderId: lookupFolderID,
		Family:   familyName,
	})
	if err != nil {
		return "", err
	}
	return image.Id, nil
}

func (c *YCClient) getInstanceIPAddress(d *Driver, instance *compute.Instance) (address string, err error) {
	// Instance could have several network interfaces with different configuration each
	// Get all possible addresses for instance
	addrIPV4Internal, addrIPV4External, addrIPV6Addr, err := c.instanceAddresses(instance)
	if err != nil {
		return "", err
	}

	if d.UseIPv6 {
		if addrIPV6Addr != "" {
			return "[" + addrIPV6Addr + "]", nil
		}
		return "", errors.New("instance has no one IPv6 address")
	}

	if d.UseInternalIP {
		if addrIPV4Internal != "" {
			return addrIPV4Internal, nil
		}
		return "", errors.New("instance has no one IPv4 internal address")
	}
	if addrIPV4External != "" {
		return addrIPV4External, nil
	}
	return "", errors.New("instance has no one IPv4 external address")
}

func (c *YCClient) instanceAddresses(instance *compute.Instance) (ipV4Int, ipV4Ext, ipV6 string, err error) {
	if len(instance.NetworkInterfaces) == 0 {
		return "", "", "", errors.New("No one network interface found for an instance")
	}

	var ipV4IntFound, ipV4ExtFound, ipV6Found bool
	for _, iface := range instance.NetworkInterfaces {
		if !ipV6Found && iface.PrimaryV6Address != nil {
			ipV6 = iface.PrimaryV6Address.Address
			ipV6Found = true
		}

		if !ipV4IntFound && iface.PrimaryV4Address != nil {
			ipV4Int = iface.PrimaryV4Address.Address
			ipV4IntFound = true

			if !ipV4ExtFound && iface.PrimaryV4Address.OneToOneNat != nil {
				ipV4Ext = iface.PrimaryV4Address.OneToOneNat.Address
				ipV4ExtFound = true
			}
		}

		if ipV6Found && ipV4IntFound && ipV4ExtFound {
			break
		}
	}

	if !ipV4IntFound {
		// internal ipV4 address always should present
		return "", "", "", errors.New("No IPv4 internal address found. Bug?")
	}

	return
}

func toBytes(gigabytesCount int) int64 {
	return int64((datasize.ByteSize(gigabytesCount) * datasize.GB).Bytes())
}
