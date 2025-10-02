//go:build service_test

package test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	"github.com/tsetsik/vehicle-search/test/setup"
	// . "github.com/onsi/gomega"
)

var _ = Describe("Search endpoint test", func() {
	var (
		suite setup.Suite
		resp  *http.Response
	)
	BeforeEach(func() {
		suite = setup.NewSuite()
	})

	Describe("GET /search", func() {
		BeforeEach(func() {
			By("Storing some vehicles")
			headers := http.Header{
				"item": []string{
					`{"make":"Toyota","model":"Camry","year":2020}`,
					`{"make":"Toyota","model":"Camry","year":2021}`,
				},
			}
			suite.MakeRequest("POST", "/listings", nil, headers)
		})

		It("finds inserted vehicles", func() {
			By("Searching for a vehicle")
			resp = suite.MakeRequest("GET", "/search?q=Toyota", nil, nil)
			fmt.Println("\n\nResponse:", resp)
		})
	})
})
