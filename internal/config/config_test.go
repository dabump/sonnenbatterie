package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	type args struct {
		fileName string
		fn       reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Successful unmarshall",
			args: args{
				fileName: "testfile.cfg",
				fn: func(name string) ([]byte, error) {
					return createValidConfigBytes(), nil
				},
			},
			wantErr: false,
			want: &Config{
				SonnenBatterieIP:                "192.168.0.0",
				SonnenBatterieStatusPath:        "/the/path/to/status",
				SonnenBatterieProtocolScheme:    "http",
				SonnenBatteriePollingTimeInMins: 3,
				HttpTimeoutInMinutes:            2,
			},
		},
		{
			name: "Unsuccessful unmarshall",
			args: args{
				fileName: "testfile.cfg",
				fn: func(name string) ([]byte, error) {
					return []byte("invalid jason"), nil
				},
			},
			wantErr: true,
		},
		{
			name: "Reader error",
			args: args{
				fileName: "testfile.cfg",
				fn: func(name string) ([]byte, error) {
					return nil, fmt.Errorf("something went wrong while reading file")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readConfigFile(tt.args.fileName, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createValidConfigBytes() []byte {
	cfg := Config{
		SonnenBatterieIP:                "192.168.0.0",
		SonnenBatterieStatusPath:        "/the/path/to/status",
		SonnenBatterieProtocolScheme:    "http",
		SonnenBatteriePollingTimeInMins: 3,
		HttpTimeoutInMinutes:            2,
	}
	bytes, _ := json.Marshal(cfg)
	return bytes
}
