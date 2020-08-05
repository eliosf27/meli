package item_test

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	itemApp "meli/internal/app/item"
	"meli/internal/entities"
	itemMocks "meli/internal/mocks/item"
	configPkg "meli/pkg/config"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("ItemHandler", func() {
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
				configs := configPkg.NewConfig()

				mockService := itemMocks.NewMockItemService(ctrl)
				mockService.EXPECT().FetchItemById(gomock.Any()).Return(entities.Item{}).AnyTimes()

				itemHandle := itemApp.NewItemHandle(configs, mockService)

				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/items/:item_id")
				c.SetParamNames("item_id")
				c.SetParamValues("")
				err := itemHandle.Get(c)

				Expect(err).To(Not(Equal(nil)))
				Expect(err.Error()).To(Equal("code=400, message=invalid item id"))
			})
		})
	})
})
