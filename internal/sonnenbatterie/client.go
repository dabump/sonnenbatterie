//go:generate mockgen -source=client.go -destination=../../test/mocks/httpclient.go -package=mocks
package sonnenbatterie

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dabump/sonnenbatterie/internal/config"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(httpClient HttpClient, config *config.Config) *Client {
	return &Client{
		config:     config,
		httpClient: httpClient,
	}
}

func (c *Client) GetStatus() (*Status, error) {

	request, _ := http.NewRequest(http.MethodGet, fmt.Sprint(c.config.SonnenBatterieProtocolScheme,
		"://", c.config.SonnenBatterieIP, c.config.SonnenBatterieStatusPath), nil)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error during invocation of status function\nstatus code:%v\nresponse: %s",
			response.StatusCode, string(bodyBytes))
	}

	status := &Status{}
	err = json.Unmarshal(bodyBytes, status)
	if err != nil {
		return nil, err
	}

	return status, nil
}
