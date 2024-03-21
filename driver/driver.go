package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/resourcemanager/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/vpc/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
)

type Driver struct {
	*drivers.BaseDriver

	Endpoint              string
	ServiceAccountKeyFile string
	Token                 string

	CloudID          string
	Cores            int
	Gpus			 int
	CoreFraction     int
	DiskSize         int
	DiskType         string
	FolderID         string
	ImageFamily      string
	ImageFolderID    string
	ImageID          string
	InstanceID       string
	Labels           []string
	Memory           int
	Metadata         map[string]string
	Nat              bool
	PlatformID       string
	Preemptible      bool
	SSHUser          string
	SubnetID         string
	UseIPv6          bool
	UseInternalIP    bool
	UserDataFile     string
	Zone             string
	StaticAddress    string
	SecurityGroups   []string
	ServiceAccountID string
	Filesystems      []string
}

const (
	defaultGpus		 = 0
	defaultCores         = 2
	defaultCoreFraction  = 100
	defaultDiskSize      = 20
	defaultDiskType      = "network-hdd"
	defaultEndpoint      = "api.cloud.yandex.net:443"
	defaultImageFamily   = "ubuntu-2004-lts"
	defaultImageFolderID = StandardImagesFolderID
	defaultMemory        = 1
	defaultPlatformID    = "standard-v1"
	defaultSSHPort       = 22
	defaultSSHUser       = "ubuntu"
	defaultZone          = "ru-central1-a"
)

func NewDriver() drivers.Driver {
	return &Driver{
		BaseDriver:    &drivers.BaseDriver{},
		Cores:         defaultCores,
		Gpus:		   defaultGpus,
		DiskSize:      defaultDiskSize,
		DiskType:      defaultDiskType,
		ImageFolderID: defaultImageFolderID,
		ImageFamily:   defaultImageFamily,
		Memory:        defaultMemory,
		Metadata:      map[string]string{},
		PlatformID:    defaultPlatformID,
		Zone:          defaultZone,
	}
}

func (d *Driver) DriverName() string {
	return "yandex"
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "YC_CLOUD_ID",
			Name:   "yandex-cloud-id",
			Usage:  "Cloud ID",
		},
		mcnflag.IntFlag{
			EnvVar: "YC_CORES",
			Name:   "yandex-cores",
			Usage:  "Count of virtual CPUs",
			Value:  defaultCores,
		},
		mcnflag.IntFlag{
			EnvVar: "YC_GPUS",
			Name:   "yandex-gpus",
			Usage:  "Count of virtual GPUs",
			Value:  defaultGpus,
		},
		mcnflag.IntFlag{
			EnvVar: "YC_CORE_FRACTION",
			Name:   "yandex-core-fraction",
			Usage:  "Core fraction",
			Value:  defaultCoreFraction,
		},
		mcnflag.IntFlag{
			EnvVar: "YC_DISK_SIZE",
			Name:   "yandex-disk-size",
			Usage:  "Disk size in gigabytes",
			Value:  defaultDiskSize,
		},
		mcnflag.StringFlag{
			EnvVar: "YC_DISK_TYPE",
			Name:   "yandex-disk-type",
			Usage:  "Disk type, e.g. 'network-hdd'",
			Value:  defaultDiskType,
		},
		mcnflag.StringFlag{
			EnvVar: "YC_ENDPOINT",
			Name:   "yandex-endpoint",
			Usage:  "Yandex.Cloud API Endpoint",
			Value:  defaultEndpoint,
		},
		mcnflag.StringFlag{
			EnvVar: "YC_FOLDER_ID",
			Name:   "yandex-folder-id",
			Usage:  "Folder ID",
		},
		mcnflag.StringFlag{
			EnvVar: "YC_IMAGE_FAMILY",
			Name:   "yandex-image-family",
			Usage:  "Image family name to lookup image ID for instance",
			Value:  defaultImageFamily,
		},
		mcnflag.StringFlag{
			EnvVar: "YC_IMAGE_FOLDER_ID",
			Name:   "yandex-image-folder-id",
			Usage:  "Folder ID to the latest image by family name",
			Value:  defaultImageFolderID,
		},
		mcnflag.StringFlag{
			EnvVar: "YC_IMAGE_ID",
			Name:   "yandex-image-id",
			Usage:  "User-defined Image ID",
			Value:  "",
		},
		mcnflag.StringSliceFlag{
			EnvVar: "YC_LABELS",
			Name:   "yandex-labels",
			Usage:  "Instance labels in 'key=value' format",
		},
		mcnflag.IntFlag{
			EnvVar: "YC_MEMORY",
			Name:   "yandex-memory",
			Usage:  "Memory in gigabytes",
			Value:  defaultMemory,
		},
		mcnflag.BoolFlag{
			EnvVar: "YC_NAT",
			Name:   "yandex-nat",
			Usage:  "Assign external (NAT) IP address",
		},
		mcnflag.StringFlag{
			EnvVar: "YC_PLATFORM_ID",
			Name:   "yandex-platform-id",
			Usage:  "ID of the hardware platform configuration",
			Value:  defaultPlatformID,
		},
		mcnflag.BoolFlag{
			EnvVar: "YC_PREEMPTIBLE",
			Name:   "yandex-preemptible",
			Usage:  "Yandex.Cloud Instance preemptibility flag",
		},
		mcnflag.StringFlag{
			EnvVar: "YC_SA_KEY_FILE",
			Name:   "yandex-sa-key-file",
			Usage:  "Yandex.Cloud Service Account key file",
		},
		mcnflag.IntFlag{
			EnvVar: "YC_SSH_PORT",
			Name:   "yandex-ssh-port",
			Usage:  "SSH port",
			Value:  defaultSSHPort,
		},
		mcnflag.StringFlag{
			EnvVar: "YC_SSH_USER",
			Name:   "yandex-ssh-user",
			Usage:  "SSH username",
			Value:  defaultSSHUser,
		},
		mcnflag.StringFlag{
			EnvVar: "YC_SUBNET_ID",
			Name:   "yandex-subnet-id",
			Usage:  "Subnet ID",
		},
		mcnflag.StringFlag{
			EnvVar: "YC_TOKEN",
			Name:   "yandex-token",
			Usage:  "Yandex.Cloud OAuth token",
		},
		mcnflag.BoolFlag{
			EnvVar: "YC_USE_INTERNAL_IP",
			Name:   "yandex-use-internal-ip",
			Usage:  "Use the internal Instance IP to communicate",
		},
		mcnflag.StringFlag{
			EnvVar: "YC_USERDATA",
			Name:   "yandex-userdata",
			Usage:  "Path to file with cloud-init user-data",
		},
		mcnflag.StringFlag{
			EnvVar: "YC_ZONE",
			Name:   "yandex-zone",
			Usage:  "Yandex.Cloud zone",
			Value:  defaultZone,
		},
		mcnflag.StringFlag{
			EnvVar: "YC_STATIC_ADDRESS",
			Name:   "yandex-static-address",
			Usage:  "Set static address",
			Value:  "",
		},
		mcnflag.StringSliceFlag{
			EnvVar: "YC_SECURITY_GROUPS",
			Name:   "yandex-security-groups",
			Usage:  "Set security groups",
		},
		mcnflag.StringFlag{
			EnvVar: "YC_SA_ID",
			Name:   "yandex-sa-id",
			Usage:  "Service account ID to attach to the instance",
		},
		mcnflag.StringSliceFlag{
			EnvVar: "YC_FS",
			Name:   "yandex-fs",
			Usage:  "Filesystem to attach to the instance. Format 'deviceName=FilesystemID'",
		},
	}
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.CloudID = flags.String("yandex-cloud-id")
	d.FolderID = flags.String("yandex-folder-id")

	d.ServiceAccountKeyFile = flags.String("yandex-sa-key-file")
	d.Token = flags.String("yandex-token")

	d.Cores = flags.Int("yandex-cores")
	d.Gpus = flags.Int("yandex-gpus")
	d.CoreFraction = flags.Int("yandex-core-fraction")
	d.DiskSize = flags.Int("yandex-disk-size")
	d.DiskType = flags.String("yandex-disk-type")
	d.Endpoint = flags.String("yandex-endpoint")
	d.ImageFamily = flags.String("yandex-image-family")
	d.ImageFolderID = flags.String("yandex-image-folder-id")
	d.ImageID = flags.String("yandex-image-id")
	d.Labels = flags.StringSlice("yandex-labels")
	d.Memory = flags.Int("yandex-memory")
	d.Nat = flags.Bool("yandex-nat")
	d.PlatformID = flags.String("yandex-platform-id")
	d.Preemptible = flags.Bool("yandex-preemptible")
	d.SSHUser = flags.String("yandex-ssh-user")
	d.SSHPort = flags.Int("yandex-ssh-port")
	d.SubnetID = flags.String("yandex-subnet-id")
	d.UseInternalIP = flags.Bool("yandex-use-internal-ip")
	d.UserDataFile = flags.String("yandex-userdata")
	d.Zone = flags.String("yandex-zone")
	d.StaticAddress = flags.String("yandex-static-address")
	d.SecurityGroups = flags.StringSlice("yandex-security-groups")
	d.ServiceAccountID = flags.String("yandex-sa-id")
	d.Filesystems = flags.StringSlice("yandex-fs")

	return nil
}

func (d *Driver) PreCreateCheck() error {
	if d.UserDataFile != "" {
		if _, err := os.Stat(d.UserDataFile); os.IsNotExist(err) {
			return fmt.Errorf("user-data file %s could not be found", d.UserDataFile)
		}
	}

	c, err := d.buildClient()
	if err != nil {
		return err
	}

	if d.FolderID == "" {
		if d.CloudID == "" {
			log.Warn("No Folder and Cloud identifiers provided")
			log.Warn("Try guess cloud ID to use")
			d.CloudID, err = d.guessCloudID()
			if err != nil {
				return err
			}
		}

		log.Warnf("Try guess folder ID to use inside cloud %q", d.CloudID)
		d.FolderID, err = d.guessFolderID()
		if err != nil {
			return err
		}
	}
	log.Infof("Check folder exists")
	folder, err := c.sdk.ResourceManager().Folder().Get(context.Background(), &resourcemanager.GetFolderRequest{
		FolderId: d.FolderID,
	})
	if err != nil {
		return fmt.Errorf("Folder with ID %q not found. %v", d.FolderID, err)
	}

	log.Infof("Check if the instance with name %q already exists in folder", d.MachineName)
	resp, err := c.sdk.Compute().Instance().List(context.Background(), &compute.ListInstancesRequest{
		FolderId: d.FolderID,
		Filter:   fmt.Sprintf("name = \"%s\"", d.MachineName),
	})
	if err != nil {
		return fmt.Errorf("Fail to get instance list in Folder: %s", err)
	}
	if len(resp.Instances) > 0 {
		return fmt.Errorf("instance with name %q already exists in folder %q", d.MachineName, d.FolderID)
	}

	if d.SubnetID == "" {
		log.Warnf("Subnet ID not provided, will search one for Zone %q in folder %q [%s]", d.Zone, folder.Name, folder.Id)
		d.SubnetID, err = d.findSubnetID()
		if err != nil {
			return err
		}

	}

	return nil
}

// Create creates a Yandex.Cloud VM instance acting as a docker host.
func (d *Driver) Create() error {
	log.Infof("Generating SSH Key")
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}

	publicKey, err := os.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return err
	}

	log.Infof("Prepare an instance metadata (user-data included)")
	if err := d.prepareInstanceMetadata(string(publicKey)); err != nil {
		return err
	}
	log.Debugf("Formed user-data:\n%s\n", d.Metadata["user-data"])

	log.Infof("Creating instance...")
	c, err := d.buildClient()
	if err != nil {
		return err
	}

	if err := c.createInstance(d); err != nil {
		// cleanup partially created instance
		_ = d.Remove()
		return err
	}

	return nil
}

// GetSSHHostname returns hostname for use with ssh
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetSSHUsername returns username for use with ssh
func (d *Driver) GetSSHUsername() string {
	if d.SSHUser == "" {
		d.SSHUser = "docker-user"
	}
	return d.SSHUser
}

// GetURL returns the URL of the remote docker daemon.
func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

func (d *Driver) GetState() (state.State, error) {
	c, err := d.buildClient()
	if err != nil {
		return state.None, err
	}

	instance, err := c.sdk.Compute().Instance().Get(context.Background(), &compute.GetInstanceRequest{
		InstanceId: d.InstanceID,
	})
	if err != nil {
		return state.Error, err
	}

	status := instance.Status
	log.Debugf("Instance State: %s", status)

	switch status {
	case compute.Instance_PROVISIONING, compute.Instance_STARTING:
		return state.Starting, nil
	case compute.Instance_RUNNING:
		return state.Running, nil
	case compute.Instance_STOPPING, compute.Instance_STOPPED, compute.Instance_DELETING:
		return state.Stopped, nil
	}

	return state.None, nil
}

func (d *Driver) Kill() error {
	return d.Stop()
}

func (d *Driver) Remove() error {
	c, err := d.buildClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	op, err := c.sdk.WrapOperation(c.sdk.Compute().Instance().Delete(ctx, &compute.DeleteInstanceRequest{
		InstanceId: d.InstanceID,
	}))
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func (d *Driver) Restart() error {
	c, err := d.buildClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	op, err := c.sdk.WrapOperation(c.sdk.Compute().Instance().Restart(ctx, &compute.RestartInstanceRequest{
		InstanceId: d.InstanceID,
	}))
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func (d *Driver) Start() error {
	c, err := d.buildClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	op, err := c.sdk.WrapOperation(c.sdk.Compute().Instance().Start(ctx, &compute.StartInstanceRequest{
		InstanceId: d.InstanceID,
	}))
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func (d *Driver) Stop() error {
	c, err := d.buildClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	op, err := c.sdk.WrapOperation(c.sdk.Compute().Instance().Stop(ctx, &compute.StopInstanceRequest{
		InstanceId: d.InstanceID,
	}))
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func (d *Driver) buildClient() (*YCClient, error) {
	return NewYCClient(d)
}

func (d *Driver) guessCloudID() (string, error) {
	c, err := d.buildClient()
	if err != nil {
		return "", err
	}

	resp, err := c.sdk.ResourceManager().Cloud().List(context.Background(), &resourcemanager.ListCloudsRequest{})
	if err != nil {
		return "", err
	}
	switch {
	case len(resp.Clouds) == 0:
		return "", errors.New("no one Cloud available")
	case len(resp.Clouds) > 1:
		return "", errors.New("more than one Cloud available, could not choose one; try specify exact Folder ID by '--yandex-folder-id' param")
	}
	return resp.Clouds[0].Id, nil
}

func (d *Driver) guessFolderID() (string, error) {
	c, err := d.buildClient()
	if err != nil {
		return "", err
	}

	resp, err := c.sdk.ResourceManager().Folder().List(context.Background(), &resourcemanager.ListFoldersRequest{
		CloudId: d.CloudID,
	})
	if err != nil {
		return "", err
	}
	switch {
	case len(resp.Folders) == 0:
		return "", errors.New("no one Folder available")
	case len(resp.Folders) > 1:
		return "", errors.New("more than one Folder available, could not choose one; try specify exact Folder ID by '--yandex-folder-id' param")
	}

	return resp.Folders[0].Id, nil
}

func (d *Driver) findSubnetID() (string, error) {
	c, err := d.buildClient()
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	resp, err := c.sdk.VPC().Subnet().List(ctx, &vpc.ListSubnetsRequest{
		FolderId: d.FolderID,
	})
	if err != nil {
		return "", err
	}

	for _, subnet := range resp.Subnets {
		if subnet.ZoneId != d.Zone {
			continue
		}
		return subnet.Id, nil
	}
	return "", fmt.Errorf("no subnets in zone: %s", d.Zone)
}

func (d *Driver) ParsedLabels() map[string]string {
	var labels = make(map[string]string)

	for _, labelPair := range d.Labels {
		labelPair = strings.TrimSpace(labelPair)
		chunks := strings.SplitN(labelPair, "=", 2)
		if len(chunks) == 1 {
			labels[chunks[0]] = ""
		} else {
			labels[chunks[0]] = chunks[1]
		}
	}
	return labels
}

func (d *Driver) ParseFilesystems() (map[string]map[string]string, error) {
	var filesystems = make(map[string]map[string]string)
	for _, fsPair := range d.Filesystems {
		fsPair = strings.TrimSpace(fsPair)
		chunks := strings.SplitN(fsPair, "=", 2)
		if len(chunks) < 2 {
			return filesystems, fmt.Errorf("wrong filesystem flag format. Need use format mountPath=FilesystemID. Example: --yandex-fs='/mnt/nfs=fs_id'")
		}
		fsPathList := strings.Split(chunks[0], "/")
		fsName := fsPathList[len(fsPathList)-1]
		filesystems[fsName] = map[string]string{"filesystemId": chunks[1], "filesystemPath": chunks[0]}
	}
	return filesystems, nil
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

func (d *Driver) prepareInstanceMetadata(publicKey string) error {
	// form 'ssh-keys' metadata key
	sshMetaDataKey := "ssh-keys"
	sshMetaDataValue := fmt.Sprintf("%s:%s", d.GetSSHUsername(), publicKey)

	d.Metadata[sshMetaDataKey] = sshMetaDataValue

	// form 'user-data' metadata key
	userData, err := d.prepareUserData(publicKey)
	if err != nil {
		return err
	}

	if userData != "" {
		d.Metadata["user-data"] = userData
	}

	return nil
}

func (d *Driver) prepareUserData(publicKey string) (string, error) {
	filesystems, _ := d.ParseFilesystems()
	userData, err := defaultUserData(d.GetSSHUsername(), publicKey, filesystems)
	if err != nil {
		return "", err
	}

	if d.UserDataFile != "" {
		log.Infof("Use provided file %q with user-data", d.UserDataFile)
		buf, err := os.ReadFile(d.UserDataFile)
		if err != nil {
			return "", err
		}
		userData, err = combineTwoCloudConfigs(userData, string(buf))
		if err != nil {
			return "", err
		}
	}

	return userData, nil
}

func (d *Driver) Credentials() (ycsdk.Credentials, error) {
	if d.ServiceAccountKeyFile != "" && d.Token != "" {
		return nil, fmt.Errorf("only one of 'token' or 'sa-key-file' should be specified")
	}

	if d.ServiceAccountKeyFile != "" {
		key, err := iamkey.ReadFromJSONFile(d.ServiceAccountKeyFile)
		if err != nil {
			return nil, err
		}
		return ycsdk.ServiceAccountKey(key)
	}

	if d.Token != "" {
		if strings.HasPrefix(d.Token, "t1.") && strings.Count(d.Token, ".") == 2 {
			return ycsdk.NewIAMTokenCredentials(d.Token), nil
		}
		return ycsdk.OAuthToken(d.Token), nil
	}

	if sa := ycsdk.InstanceServiceAccount(); checkServiceAccountAvailable(context.Background(), sa) {
		return sa, nil
	}

	return nil, fmt.Errorf("one of 'token' or 'sa-key-file' should be specified; if you are inside compute instance, you can attach service account to it in order to authenticate via instance service account")
}

func checkServiceAccountAvailable(ctx context.Context, sa ycsdk.NonExchangeableCredentials) bool {
	dialer := net.Dialer{Timeout: 50 * time.Millisecond}
	conn, err := dialer.Dial("tcp", net.JoinHostPort(ycsdk.InstanceMetadataAddr, "80"))
	if err != nil {
		return false
	}
	_ = conn.Close()
	_, err = sa.IAMToken(ctx)
	return err == nil
}

func defaultUserData(sshUserName, sshPublicKey string, fs map[string]map[string]string) (string, error) {
	type templateData struct {
		SSHUserName  string
		SSHPublicKey string
		Filesystems  map[string]map[string]string
	}
	buf := &bytes.Buffer{}
	err := defaultUserDataTemplate.Execute(buf, templateData{
		SSHUserName:  sshUserName,
		SSHPublicKey: sshPublicKey,
		Filesystems:  fs,
	})
	if err != nil {
		return "", fmt.Errorf("error while process template: %s", err)
	}

	return buf.String(), nil
}

var defaultUserDataTemplate = template.Must(
	template.New("user-data").Parse(`#cloud-config
ssh_pwauth: no

users:
  - name: {{.SSHUserName}}
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - {{.SSHPublicKey}}

{{ if gt (len .Filesystems) 0}}
runcmd:
{{range $name, $fs := .Filesystems}}
  - mkdir {{index $fs "filesystemPath"}}
  - mount -t virtiofs {{ $name }} {{index $fs "filesystemPath"}}
{{end}}
{{end}}
`))
