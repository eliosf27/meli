package api

import (
	"github.com/dghubble/sling"
	"net/http"
	"time"
)

type HttpClient struct {
	ItemService *ItemService
}

func NewHttpClient() HttpClient {
	client := http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	sl := sling.New().Client(&client).Base("https://api.mercadolibre.com/")

	return HttpClient{
		ItemService: NewItemService(sl),
	}
}
