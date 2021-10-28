package v1_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/epinio/epinio/acceptance/helpers/catalog"
	api "github.com/epinio/epinio/internal/api/v1"
	"github.com/epinio/epinio/pkg/api/core/v1/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services API Application Endpoints", func() {
	containerImageURL := "splatform/sample-app"

	var org string
	var svc1, svc2 string

	BeforeEach(func() {
		org = catalog.NewOrgName()
		env.SetupAndTargetOrg(org)

		svc1 = catalog.NewServiceName()
		svc2 = catalog.NewServiceName()

		env.MakeService(svc1)
		env.MakeService(svc2)
	})

	AfterEach(func() {
		env.DeleteService(svc1)
		env.DeleteService(svc2)
	})

	Describe("GET /api/v1/namespaces/:org/services", func() {
		var serviceNames []string

		It("lists all services in the org", func() {
			response, err := env.Curl("GET", fmt.Sprintf("%s%s/namespaces/%s/services",
				serverURL, api.Root, org), strings.NewReader(""))
			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())
			defer response.Body.Close()
			bodyBytes, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK), string(bodyBytes))

			var data models.ServiceResponseList
			err = json.Unmarshal(bodyBytes, &data)
			Expect(err).ToNot(HaveOccurred())
			serviceNames = append(serviceNames, data[0].Meta.Name)
			serviceNames = append(serviceNames, data[1].Meta.Name)
			Expect(serviceNames).Should(ContainElements(svc1, svc2))
		})

		It("returns a 404 when the org does not exist", func() {
			response, err := env.Curl("GET", fmt.Sprintf("%s%s/namespaces/idontexist/services",
				serverURL, api.Root),
				strings.NewReader(""))
			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())

			defer response.Body.Close()
			bodyBytes, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusNotFound), string(bodyBytes))
		})
	})

	Describe("GET /api/v1/services", func() {
		var org1 string
		var org2 string
		var service1 string
		var service2 string
		var app1 string

		// Setting up:
		// org1 service1 app1
		// org2 service1
		// org2 service2

		BeforeEach(func() {
			org1 = catalog.NewOrgName()
			org2 = catalog.NewOrgName()
			service1 = catalog.NewServiceName()
			service2 = catalog.NewServiceName()
			app1 = catalog.NewAppName()

			env.SetupAndTargetOrg(org1)
			env.MakeService(service1)
			env.MakeContainerImageApp(app1, 1, containerImageURL)
			env.BindAppService(app1, service1, org1)

			env.SetupAndTargetOrg(org2)
			env.MakeService(service1) // separate from org1.service1
			env.MakeService(service2)
		})

		AfterEach(func() {
			env.TargetOrg(org2)
			env.DeleteService(service1)
			env.DeleteService(service2)

			env.TargetOrg(org1)
			env.DeleteApp(app1)
			env.DeleteService(service1)
		})

		It("lists all services belonging to all namespaces", func() {
			// But we care only about the three we know about from the setup.

			response, err := env.Curl("GET", fmt.Sprintf("%s%s/services",
				serverURL, api.Root), strings.NewReader(""))
			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())

			defer response.Body.Close()
			bodyBytes, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK), string(bodyBytes))

			var services models.ServiceResponseList
			err = json.Unmarshal(bodyBytes, &services)
			Expect(err).ToNot(HaveOccurred())

			// `services` contains all services. Not just the two we are looking for, from
			// the setup of this test. Everything which still exists from other tests
			// executing concurrently, or not cleaned by previous tests, or the setup, or
			// ... So, we cannot be sure that the two services are in the two first
			// elements of the slice.

			var serviceRefs [][]string
			for _, s := range services {
				serviceRefs = append(serviceRefs, []string{
					s.Meta.Name,
					s.Meta.Namespace,
					strings.Join(s.BoundApps, ", "),
				})
			}
			Expect(serviceRefs).To(ContainElements(
				[]string{service1, org1, app1},
				[]string{service1, org2, ""},
				[]string{service2, org2, ""},
			))

		})
	})

	Describe("GET /api/v1/namespaces/:org/services/:service", func() {
		It("lists the service data", func() {
			response, err := env.Curl("GET", fmt.Sprintf("%s%s/namespaces/%s/services/%s",
				serverURL, api.Root, org, svc1), strings.NewReader(""))
			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())
			defer response.Body.Close()
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			bodyBytes, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())

			var data models.ServiceShowResponse
			err = json.Unmarshal(bodyBytes, &data)
			service := data.Details
			Expect(err).ToNot(HaveOccurred())
			Expect(service["username"]).To(Equal("epinio-user"))
		})

		It("returns a 404 when the org does not exist", func() {
			response, err := env.Curl("GET", fmt.Sprintf("%s%s/namespaces/idontexist/services/%s",
				serverURL, api.Root, svc1), strings.NewReader(""))
			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())

			defer response.Body.Close()
			bodyBytes, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusNotFound), string(bodyBytes))
		})

		It("returns a 404 when the service does not exist", func() {
			response, err := env.Curl("GET", fmt.Sprintf("%s%s/namespaces/%s/services/bogus",
				serverURL, api.Root, org), strings.NewReader(""))
			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())

			defer response.Body.Close()
			bodyBytes, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusNotFound), string(bodyBytes))
		})
	})
})
