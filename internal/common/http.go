package common

import "net/http"

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=http.go -destination=../../test/mocks/httpclient.go -package=mocks
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func InternalServerError(resp http.ResponseWriter, err error) {
	resp.Header().Set("Content-Type", "application/text")
	resp.WriteHeader(http.StatusInternalServerError)
	_, _ = resp.Write([]byte(err.Error()))
}

func TooManyRequests(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/text")
	resp.WriteHeader(http.StatusTooManyRequests)
}
