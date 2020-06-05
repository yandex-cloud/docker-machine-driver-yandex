package driver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDriver_prepareInstanceMetadata(t *testing.T) {
	mockSshPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDkai1XE7djYB5Z"

	type fields struct {
		SSHUser      string
		UserDataFile string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantMD  map[string]string
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
			},
		},
		{
			name: "user-data from file",
			fields: fields{
				SSHUser:      "debian",
				UserDataFile: "test-fixtures/user-data.txt",
			},
			wantErr: false,
			wantMD: map[string]string{
				"ssh-keys":  "debian:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDkai1XE7djYB5Z",
				"user-data": "My Custom User-Data",
			},
		},
		{
			name: "user-data file does not exist",
			fields: fields{
				SSHUser:      "debian",
				UserDataFile: "test-fixtures/no-such-file.txt",
			},
			wantErr: true,
			wantMD:  map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Driver{
				Metadata:     map[string]string{},
				SSHUser:      tt.fields.SSHUser,
				UserDataFile: tt.fields.UserDataFile,
			}
			e := d.prepareInstanceMetadata(mockSshPublicKey)
			if tt.wantErr {
				require.Error(t, e, "error expected")
			} else {
				require.NoError(t, e, "no error expected, got one")
				require.Equal(t, tt.wantMD, d.Metadata)
			}
		})
	}
}
