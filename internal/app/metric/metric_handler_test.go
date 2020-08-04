package metric

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"meli/internal/entities"
	configMetric "meli/internal/mocks/metric"
	configPkg "meli/pkg/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var _ = Describe("MetricHandler", func() {
	var ctrl *gomock.Controller

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Context("calling the fetch metric handler", func() {
		When("the request is valid", func() {
			It("should return a valid response", func() {
				configs := configPkg.NewConfig()

				response := []entities.ItemMetricResponse{
					{
						AvgResponseTime:        15.8,
						TotalRequests:          24,
						AvgResponseTimeApiCall: 58.1,
						TotalCountApiCalls:     15,
						Time:                   time.Now(),
						InfoRequests:           nil,
					},
				}

				mockService := configMetric.NewMockMetricService(ctrl)
				mockService.EXPECT().FetchMetrics().Return(response).AnyTimes()

				metricHandler := NewMetricHandler(configs, mockService)

				data, err := json.Marshal(response)
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(string(data)))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("health")

				err = metricHandler.Health(c)
				Expect(err).To(BeNil())

				logrus.Info(rec.Body.String())
				logrus.Info(string(data))

				var expectedResponse []entities.ItemMetricResponse
				err = json.Unmarshal([]byte(rec.Body.String()), &expectedResponse)
				Expect(err).To(BeNil())
				Expect(len(expectedResponse)).To(Equal(len(response)))
				Expect(expectedResponse[0].AvgResponseTimeApiCall).To(Equal(response[0].AvgResponseTimeApiCall))
				Expect(expectedResponse[0].AvgResponseTime).To(Equal(response[0].AvgResponseTime))
				Expect(expectedResponse[0].TotalCountApiCalls).To(Equal(response[0].TotalCountApiCalls))
				Expect(expectedResponse[0].TotalRequests).To(Equal(response[0].TotalRequests))
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
