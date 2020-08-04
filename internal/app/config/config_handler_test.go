package config

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/go-playground/validator.v9"
	configMocks "meli/internal/mocks/config"
	"net/http"
	"net/http/httptest"
	"strings"
)

type APIValidator struct {
	validator *validator.Validate
}

func (cv *APIValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var _ = Describe("ConfigHandler", func() {
	var ctrl *gomock.Controller

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Context("calling the fetch config handler", func() {
		When("the request is valid", func() {
			It("should return a valid response", func() {
				mockService := configMocks.NewMockConfigServicer(ctrl)
				mockService.EXPECT().Fetch().Return("redis", nil).AnyTimes()

				configCacheHandle := NewConfigCacheHandle(mockService)

				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("config/cache")
				err := configCacheHandle.Fetch(c)

				Expect(err).To(BeNil())
				Expect(rec.Body.String()).To(Equal("\"redis\"\n"))
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("calling the update config handler", func() {
		When("the request is valid", func() {
			It("should return a valid response", func() {
				returnStorage := "redis"
				mockService := configMocks.NewMockConfigServicer(ctrl)
				mockService.EXPECT().Update(returnStorage).Return(nil).AnyTimes()

				configCacheHandle := NewConfigCacheHandle(mockService)

				userJSON := `{"storage": "redis"}`
				e := echo.New()
				e.Validator = &APIValidator{validator: validator.New()}
				req := httptest.NewRequest(http.MethodPut, "/config/cache/", strings.NewReader(userJSON))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				err := configCacheHandle.Update(c)
				Expect(err).To(BeNil())
				Expect(rec.Body.String()).To(Equal("\"storage updated\"\n"))
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})

		When("the request is invalid", func() {
			It("should return an error", func() {
				returnStorage := "redis"
				mockService := configMocks.NewMockConfigServicer(ctrl)
				mockService.EXPECT().Update(returnStorage).Return(nil).AnyTimes()

				configCacheHandle := NewConfigCacheHandle(mockService)

				userJSON := `{"storage": ""}`
				e := echo.New()
				e.Validator = &APIValidator{validator: validator.New()}
				req := httptest.NewRequest(http.MethodPut, "/config/cache/", strings.NewReader(userJSON))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				err := configCacheHandle.Update(c)
				Expect(err).To(Not(BeNil()))
				Expect(rec.Body.String()).To(Equal(""))
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
