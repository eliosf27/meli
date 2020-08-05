package item_test

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appConfig "meli/internal/app/config"
	appItem "meli/internal/app/item"
	"meli/internal/entities"
	"meli/internal/http"
	itemMocks "meli/internal/mocks/item"
	"meli/internal/postgres"
	"meli/internal/queue"
	redis2 "meli/internal/redis"
	config "meli/pkg/config"
	mocksPkg "meli/pkg/mocks"
)

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
			It("should return an valid item", func() {
				configs := config.NewConfig()
				redis := redis2.NewRedis(configs)
				postgres := postgres.NewPostgres(configs)
				itemQueue := queue.NewItemQueue()
				httpClient := http.NewHttpClient(configs, &itemQueue)
				itemHttpService := http.NewItemHttpService(&httpClient)
				configService := appConfig.NewConfigService(redis)
				itemPostgresCache := appItem.NewItemPostgresCache(postgres)

				itemId := "jjj"
				anItem := entities.Item{}
				mockItemCacher := itemMocks.NewMockItemCacher(ctrl)
				mockItemCacher.EXPECT().Save(entities.Item{ItemId: itemId}).AnyTimes()
				mockItemCacher.EXPECT().Get(gomock.Any()).Return(anItem, errors.New("not valid item")).AnyTimes()
				mockItemCacher.EXPECT().GetStrategy(gomock.Any()).Return(itemPostgresCache).AnyTimes()
				mockItemCacher.EXPECT().SetStrategy(itemPostgresCache).AnyTimes()

				itemService := appItem.NewItemService(itemHttpService, mockItemCacher, configService)

				item := itemService.FetchItemById(itemId)

				Expect(item.IsZero()).To(Equal(false))
			})
		})

		When("a valid item is returned by the repository", func() {
			It("should return a not empty item", func() {
				itemId := "yyy"
				configs := config.NewConfig()
				redis := redis2.NewRedis(configs)
				postgres := postgres.NewPostgres(configs)
				itemQueue := queue.NewItemQueue()
				httpClient := http.NewHttpClient(configs, &itemQueue)
				itemHttpService := http.NewItemHttpService(&httpClient)
				configService := appConfig.NewConfigService(redis)
				itemPostgresCache := appItem.NewItemPostgresCache(postgres)

				mockItemCacher := itemMocks.NewMockItemCacher(ctrl)
				mockItemCacher.EXPECT().Get(gomock.Any()).Return(entities.Item{ItemId: itemId}, nil).AnyTimes()
				mockItemCacher.EXPECT().GetStrategy(gomock.Any()).Return(itemPostgresCache).AnyTimes()
				mockItemCacher.EXPECT().SetStrategy(itemPostgresCache).AnyTimes()

				itemService := appItem.NewItemService(itemHttpService, mockItemCacher, configService)

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
				redis := redis2.NewRedis(configs)
				postgres := postgres.NewPostgres(configs)
				itemQueue := queue.NewItemQueue()
				httpClient := http.NewHttpClient(configs, &itemQueue)
				itemHttpService := http.NewItemHttpService(&httpClient)
				configService := appConfig.NewConfigService(redis)
				itemPostgresCache := appItem.NewItemPostgresCache(postgres)

				// httpserver itemMocks
				itemChildrenPath := fmt.Sprintf(
					"%s%s", configs.BaseEndpoint, itemHttpService.GetItemChildrenPath(itemId),
				)
				itemPath := fmt.Sprintf(
					"%s%s", configs.BaseEndpoint, itemHttpService.GetItemPath(itemId),
				)
				httpMock.Get(itemPath, itemMocks.MockItem(itemId))
				httpMock.Get(itemChildrenPath, itemMocks.MockItemChildren(itemId))

				// repository itemMocks
				mockItemCacher := itemMocks.NewMockItemCacher(ctrl)
				mockItemCacher.EXPECT().Get(gomock.Any()).Return(entities.Item{ItemId: ""}, nil).AnyTimes()
				mockItemCacher.EXPECT().Save(gomock.Any()).Return(nil).AnyTimes()
				mockItemCacher.EXPECT().SetStrategy(itemPostgresCache).AnyTimes()
				mockItemCacher.EXPECT().GetStrategy(gomock.Any()).Return(itemPostgresCache).AnyTimes()
				//mockItemCacher.SetStrategy(appConfig.PostgresStorage)

				// ItemService to test
				itemService := appItem.NewItemService(itemHttpService, mockItemCacher, configService)
				item := itemService.FetchItemById(itemId)

				Expect(item).To(Not(Equal(item.IsZero())))
				Expect(itemId).To(Equal(item.ItemId))
			})
		})
	})
})
