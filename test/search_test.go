//go:build service_test

package test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tsetsik/vehicle-search/internal/core"
	"github.com/tsetsik/vehicle-search/test/setup"
)

type (
	SearchResponse struct {
		Items []core.Item `json:"items"`
	}
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
					`{"make":"Toyota","model":"Camry","year":2020, "description": "A reliable sedan", "fuel_type": "diesel"}`,
				},
			}
			suite.MakeRequest("POST", "/listings", nil, headers)
			headers = http.Header{
				"item": []string{
					`{"make":"Audi","model":"RS3","year":2020, "description": "A sporty sedan", "fuel_type": "petrol"}`,
				},
			}
			suite.MakeRequest("POST", "/listings", nil, headers)
		})

		It("finds inserted vehicles", func() {
			By("Searching for a vehicle")
			headers := http.Header{
				"Content-Type": []string{"application/json"},
				"q":            []string{"sporty"},
			}
			resp = suite.MakeRequest("GET", "/search", nil, headers)

			searchResp := &SearchResponse{}
			json.NewDecoder(resp.Body).Decode(searchResp)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(len(searchResp.Items)).To(Equal(1))
			Expect(searchResp.Items[0].Model).To(Equal("RS3"))
		})
	})
})
