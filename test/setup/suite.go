//go:build service_test

package setup

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/gomega"
	"github.com/tsetsik/vehicle-search/internal/config"
	httpsvc "github.com/tsetsik/vehicle-search/internal/http"
)

type (
	Suite struct {
		ctx    context.Context
		cfg    config.Config
		client *http.Client
	}
)

func NewSuite() Suite {
	ctx := context.Background()

	service, err := httpsvc.NewService()
	Expect(err).NotTo(HaveOccurred())

	go func() {
		if err := service.Start(ctx); err != nil {
			panic(err)
		}
	}()

	// Wait for the service to start
	time.Sleep(500 * time.Millisecond)

	return Suite{
		ctx:    context.Background(),
		cfg:    service.Config(),
		client: http.DefaultClient,
	}
}

func (s *Suite) MakeRequest(method, path string, body any, headers http.Header) *http.Response {
	host := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	req := &http.Request{
		Method: method,
		URL:    &url.URL{Scheme: "http", Host: host, Path: path},
		Header: headers,
	}
	resp, err := s.client.Do(req)

	Expect(err).NotTo(HaveOccurred())

	Expect(resp.StatusCode).ToNot(Equal(http.StatusInternalServerError))
	Expect(resp.StatusCode).ToNot(Equal(http.StatusNotFound))

	return resp
}
