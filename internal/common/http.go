package common

import "net/http"

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=http.go -destination=../../test/mocks/httpclient.go -package=mocks
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
