package api

import (
	"github.com/dghubble/sling"
	config "meli/pkg/config"
	"net/http"
	"time"
)

type HttpClient struct {
	ItemService *ItemService
}

func NewHttpClient(config config.Config) HttpClient {
	client := http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	sl := sling.New().Client(&client).Base(config.BaseEndpoint)

	return HttpClient{
		ItemService: NewItemService(sl),
	}
}
