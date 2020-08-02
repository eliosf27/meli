package item_test

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"meli/internal/app/entities"
	"meli/internal/app/item"
	"meli/internal/http"
	mocks "meli/internal/mocks"
	"meli/internal/queue"
	config "meli/pkg/config"
	mocksPkg "meli/pkg/mocks"
)

var _ = Describe("ItemHttpService", func() {
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
				itemQueue := queue.NewItemQueue()
				httpClient := http.NewHttpClient(configs, &itemQueue)
				itemHttpService := http.NewItemHttpService(&httpClient)

				mockRepository := mocks.NewMockItemRepository(ctrl)
				mockRepository.EXPECT().Get(gomock.Any()).Return(entities.Item{}, errors.New("not valid item")).AnyTimes()

				itemService := item.NewItemService(itemHttpService, mockRepository)

				item := itemService.FetchItemById("xxx")

				Expect(item.IsZero()).To(Equal(true))
			})
		})

		When("a valid item is returned by the repository", func() {
			It("should return a not empty item", func() {
				itemId := "yyy"
				configs := config.NewConfig()
				itemQueue := queue.NewItemQueue()
				httpClient := http.NewHttpClient(configs, &itemQueue)
				itemHttpService := http.NewItemHttpService(&httpClient)

				mockRepository := mocks.NewMockItemRepository(ctrl)
				mockRepository.EXPECT().Get(gomock.Any()).Return(entities.Item{ItemId: itemId}, nil).AnyTimes()

				itemService := item.NewItemService(itemHttpService, mockRepository)

				item := itemService.FetchItemById(itemId)

				Expect(item).To(Not(Equal(item.IsZero())))
				Expect(itemId).To(Equal(item.ItemId))
			})
		})

		When("an empty item is returned by the repository", func() {
			It("should request the http and return a valid item", func() {
				httpMock := mocksPkg.NewHttpMock()
				httpMock.Activate()
				defer httpMock.DeactivateAndReset()

				// dependencies
				itemId := "yyy"
				configs := config.NewConfig()
				itemQueue := queue.NewItemQueue()
				httpClient := http.NewHttpClient(configs, &itemQueue)
				itemHttpService := http.NewItemHttpService(&httpClient)

				// httpserver mocks
				itemChildrenPath := fmt.Sprintf(
					"%s%s", configs.BaseEndpoint, itemHttpService.GetItemChildrenPath(itemId),
				)
				itemPath := fmt.Sprintf(
					"%s%s", configs.BaseEndpoint, itemHttpService.GetItemPath(itemId),
				)
				httpMock.Get(itemPath, mocks.MockItem(itemId))
				httpMock.Get(itemChildrenPath, mocks.MockItemChildren(itemId))

				// repository mocks
				mockRepository := mocks.NewMockItemRepository(ctrl)
				mockRepository.EXPECT().Get(gomock.Any()).Return(entities.Item{ItemId: ""}, nil).AnyTimes()
				mockRepository.EXPECT().Save(gomock.Any()).Return(nil).AnyTimes()

				// service to test
				itemService := item.NewItemService(itemHttpService, mockRepository)
				item := itemService.FetchItemById(itemId)

				Expect(item).To(Not(Equal(item.IsZero())))
				Expect(itemId).To(Equal(item.ItemId))
			})
		})
	})
})
