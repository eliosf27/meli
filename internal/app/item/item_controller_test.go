package item_test

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"meli/internal/app/entities"
	"meli/internal/app/item"
	mocks "meli/internal/mocks"
	config "meli/pkg/config"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("ItemController", func() {
	var ctrl *gomock.Controller
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Context("calling the item controller", func() {
		When("the item_is is empty", func() {
			It("should return an error response", func() {
				configs := config.NewConfig()

				mockService := mocks.NewMockItemService(ctrl)
				mockService.EXPECT().FetchItemById(gomock.Any()).Return(entities.Item{}).AnyTimes()

				itemController := item.NewItemController(configs, mockService)

				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/items/:item_id")
				c.SetParamNames("item_id")
				c.SetParamValues("")
				err := itemController.Get(c)

				Expect(err).To(Not(Equal(nil)))
				Expect(err.Error()).To(Equal("code=400, message=invalid item id"))
			})
		})

		When("the item_is is not empty", func() {
			It("should return an successful response", func() {
				configs := config.NewConfig()

				mockService := mocks.NewMockItemService(ctrl)
				mockService.EXPECT().FetchItemById(gomock.Any()).Return(entities.Item{}).AnyTimes()

				itemController := item.NewItemController(configs, mockService)

				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/items/:item_id")
				c.SetParamNames("item_id")
				c.SetParamValues("MLU460998489")
				err := itemController.Get(c)

				Expect(err).To(BeNil())
			})
		})
	})
})
