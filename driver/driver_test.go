package driver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDriver_prepareInstanceMetadata(t *testing.T) {
	mockSshPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDkai1XE7djYB5Z"

	type fields struct {
		SSHUser      string
		UserDataFile string
		Filesystems  []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantMD  map[string]string
		golden  string
	}{
		{
			name: "no user-data input",
			fields: fields{
				SSHUser:      "ubuntu",
				UserDataFile: "",
			},
			wantErr: false,
			wantMD: map[string]string{
				"ssh-keys": "ubuntu:" + mockSshPublicKey,
				"user-data": `#cloud-config
ssh_pwauth: no

users:
  - name: ubuntu
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDkai1XE7djYB5Z


`,
			},
			golden: "no_user-data_input",
		},
		{
			name: "fs-user-data",
			fields: fields{
				SSHUser:      "ubuntu",
				UserDataFile: "",
				Filesystems:  []string{"/data=qwdvj7dgfksdfd"},
			},
			wantErr: false,
			wantMD: map[string]string{
				"ssh-keys": "ubuntu:" + mockSshPublicKey,
				"user-data": `#cloud-config
ssh_pwauth: no

users:
- name: ubuntu
sudo: ALL=(ALL) NOPASSWD:ALL
shell: /bin/bash
ssh_authorized_keys:
- ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDkai1XE7djYB5Z

runcmd:

- mkdir /data
- mount -t virtiofs data /data


`,
			},
			golden: "fs-user-data",
		},
		{
			name: "user-data from file",
			fields: fields{
				SSHUser:      "debian",
				UserDataFile: "testdata/user-data.txt",
			},
			wantErr: false,
			wantMD: map[string]string{
				"ssh-keys": "debian:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDkai1XE7djYB5Z",
			},
			golden: "user-data_from_file",
		},
		{
			name: "user-data file does not exist",
			fields: fields{
				SSHUser:      "debian",
				UserDataFile: "test-fixtures/no-such-file.txt",
			},
			wantErr: true,
			wantMD:  map[string]string{},
			golden:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Driver{
				Metadata:     map[string]string{},
				SSHUser:      tt.fields.SSHUser,
				UserDataFile: tt.fields.UserDataFile,
				Filesystems:  tt.fields.Filesystems,
			}
			e := d.prepareInstanceMetadata(mockSshPublicKey)
			if tt.wantErr {
				require.Error(t, e, "error expected")
			} else {
				require.NoError(t, e, "no error expected, got one")
				content, err := os.ReadFile("testdata/" + tt.golden + ".golden")
				if err != nil {
					t.Fatalf("Error loading golden file: %s", err)
				}
				want := string(content)
				require.Equal(t, want, d.Metadata["user-data"])
			}
		})
	}
}
