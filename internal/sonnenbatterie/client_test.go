package sonnenbatterie

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/dabump/sonnenbatterie/internal/common"
	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/test/mocks"
	"github.com/golang/mock/gomock"
)

func newSentientHttpClientMock(t *testing.T) common.HttpClient {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockHttpClient(ctrl)

	// -- Successful
	successfulBody := "{\"Apparent_output\":240,\"BackupBuffer\":\"0\",\"BatteryCharging\":false,\"BatteryDischarging\":false,\"Consumption_Avg\":565,\"Consumption_W\":560,\"Fac\":50.0099983215332,\"FlowConsumptionBattery\":false,\"FlowConsumptionGrid\":false,\"FlowConsumptionProduction\":true,\"FlowGridBattery\":false,\"FlowProductionBattery\":false,\"FlowProductionGrid\":true,\"GridFeedIn_W\":1960,\"IsSystemInstalled\":1,\"OperatingMode\":\"2\",\"Pac_total_W\":-5,\"Production_W\":2525,\"RSOC\":100,\"RemainingCapacity_Wh\":8283,\"Sac1\":80,\"Sac2\":80,\"Sac3\":79,\"SystemStatus\":\"OnGrid\",\"Timestamp\":\"2023-03-01 15:02:45\",\"USOC\":100,\"Uac\":243,\"Ubat\":55,\"dischargeNotAllowed\":false,\"generator_autostart\":false}"
	successfulRequest, _ := http.NewRequest(http.MethodGet, "http://192.168.2.101/the/successful/path", nil)
	successfulReponse := http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(successfulBody)),
	}
	mockClient.EXPECT().
		Do(gomock.Eq(successfulRequest)).
		Return(&successfulReponse, nil).
		Times(1)

	// -- Uncuccessful - 500
	unSuccessful5xxBody := "{\"error\":\"error occured\"}"
	unSuccessful5xxRequest, _ := http.NewRequest(http.MethodGet, "http://192.168.2.101/the/unsuccessful-500/path", nil)
	unSuccessful5xxReponse := http.Response{
		StatusCode: 501,
		Body:       io.NopCloser(strings.NewReader(unSuccessful5xxBody)),
	}
	mockClient.EXPECT().
		Do(gomock.Eq(unSuccessful5xxRequest)).
		Return(&unSuccessful5xxReponse, nil).
		Times(1)

	// -- Uncuccessful - http client
	unSuccessfulHttpclientRequest, _ := http.NewRequest(http.MethodGet, "http://192.168.2.101/the/unsuccessful-http-client/path", nil)
	mockClient.EXPECT().
		Do(gomock.Eq(unSuccessfulHttpclientRequest)).
		Return(nil, fmt.Errorf("error during invoking http client logic")).
		Times(1)
	return mockClient
}

func TestClient_GetStatus(t *testing.T) {
	sentientHttpClientMock := newSentientHttpClientMock(t)
	type fields struct {
		httpClient common.HttpClient
		config     *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Status
		wantErr bool
	}{
		{
			name: "Unsuccessful - http client failure",
			fields: fields{
				config: &config.Config{
					SonnenBatterieIP:             "192.168.2.101",
					SonnenBatterieStatusPath:     "/the/unsuccessful-http-client/path",
					SonnenBatterieProtocolScheme: "http",
					HttpTimeoutInMinutes:         1,
				},
				httpClient: sentientHttpClientMock,
			},
			wantErr: true,
		},
		{
			name: "Unsuccessful - 5xx error",
			fields: fields{
				config: &config.Config{
					SonnenBatterieIP:             "192.168.2.101",
					SonnenBatterieStatusPath:     "/the/unsuccessful-500/path",
					SonnenBatterieProtocolScheme: "http",
					HttpTimeoutInMinutes:         1,
				},
				httpClient: sentientHttpClientMock,
			},
			wantErr: true,
		},
		{
			name: "Successful status fetch",
			fields: fields{
				config: &config.Config{
					SonnenBatterieIP:             "192.168.2.101",
					SonnenBatterieStatusPath:     "/the/successful/path",
					SonnenBatterieProtocolScheme: "http",
					HttpTimeoutInMinutes:         1,
				},
				httpClient: sentientHttpClientMock,
			},
			wantErr: false,
			want: &Status{
				ApparentOutput:            240,
				BackupBuffer:              "0",
				BatteryCharging:           false,
				BatteryDischarging:        false,
				ConsumptionAvg:            565,
				ConsumptionW:              560,
				Fac:                       50.0099983215332,
				FlowConsumptionBattery:    false,
				FlowConsumptionGrid:       false,
				FlowConsumptionProduction: true,
				FlowGridBattery:           false,
				FlowProductionBattery:     false,
				FlowProductionGrid:        true,
				GridFeedInW:               1960,
				IsSystemInstalled:         1,
				OperatingMode:             "2",
				PacTotalW:                 -5,
				ProductionW:               2525,
				Rsoc:                      100,
				RemainingCapacityWh:       8283,
				Sac1:                      80,
				Sac2:                      80,
				Sac3:                      79,
				SystemStatus:              "OnGrid",
				Timestamp:                 "2023-03-01 15:02:45",
				Usoc:                      100,
				Uac:                       243,
				Ubat:                      55,
				DischargeNotAllowed:       false,
				GeneratorAutostart:        false,
			},
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient: tt.fields.httpClient,
				config:     tt.fields.config,
			}
			got, err := c.GetStatus(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
