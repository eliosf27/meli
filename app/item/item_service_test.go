package item_test

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"meli/app/entities"
	"meli/app/item"
	"meli/internal/api"
	mocks "meli/internal/mocks"
	config "meli/pkg/config"
	mocksPkg "meli/pkg/mocks"
	"testing"
)

func TestIte3m(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Item Suite")
}

var _ = Describe("ItemService", func() {
	var ctrl *gomock.Controller
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Context("calling the FetchItemById method", func() {
		When("an error is returned by the repository", func() {
			It("should return an empty item", func() {
				configs := config.NewConfig()
				httpClient := api.NewHttpClient(configs)

				mockRepository := mocks.NewMockItemRepository(ctrl)
				mockRepository.EXPECT().Get(gomock.Any()).Return(entities.Item{}, errors.New("not valid item")).AnyTimes()

				itemService := item.NewItemService(httpClient, mockRepository)

				item := itemService.FetchItemById("xxx")

				Expect(item.IsZero()).To(Equal(true))
			})
		})

		When("a valid item is returned by the repository", func() {
			It("should return a not empty item", func() {
				itemId := "yyy"
				configs := config.NewConfig()
				httpClient := api.NewHttpClient(configs)

				mockRepository := mocks.NewMockItemRepository(ctrl)
				mockRepository.EXPECT().Get(gomock.Any()).Return(entities.Item{ItemId: itemId}, nil).AnyTimes()

				itemService := item.NewItemService(httpClient, mockRepository)

				item := itemService.FetchItemById(itemId)

				Expect(item).To(Not(Equal(item.IsZero())))
				Expect(itemId).To(Equal(item.ItemId))
			})
		})

		When("an empty item is returned by the repository", func() {
			It("should request the api and return a valid item", func() {
				httpMock := mocksPkg.NewHttpMock()
				httpMock.Activate()
				defer httpMock.DeactivateAndReset()

				// dependencies
				itemId := "yyy"
				configs := config.NewConfig()
				httpClient := api.NewHttpClient(configs)

				// http mocks
				itemChildrenPath := fmt.Sprintf(
					"%s%s", configs.BaseEndpoint, httpClient.ItemService.GetItemChildrenPath(itemId),
				)
				itemPath := fmt.Sprintf(
					"%s%s", configs.BaseEndpoint, httpClient.ItemService.GetItemPath(itemId),
				)
				httpMock.Get(itemPath, mocks.MockItem(itemId))
				httpMock.Get(itemChildrenPath, mocks.MockItemChildren(itemId))

				// repository mocks
				mockRepository := mocks.NewMockItemRepository(ctrl)
				mockRepository.EXPECT().Get(gomock.Any()).Return(entities.Item{ItemId: ""}, nil).AnyTimes()
				mockRepository.EXPECT().Save(gomock.Any()).Return(nil).AnyTimes()

				// service to test
				itemService := item.NewItemService(httpClient, mockRepository)
				item := itemService.FetchItemById(itemId)

				Expect(item).To(Not(Equal(item.IsZero())))
				Expect(itemId).To(Equal(item.ItemId))
			})
		})
	})
})
