package http

import (
	"github.com/dghubble/sling"
	"meli/internal/app/entities"
	"meli/internal/queue"
	config "meli/pkg/config"
	"net/http"
	"time"
)

type HttpClient struct {
	sling *sling.Sling
	queue *queue.ItemQueue
}

func NewHttpClient(config config.Config, queue *queue.ItemQueue) HttpClient {
	client := http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	return HttpClient{
		sling: sling.New().Client(&client).Base(config.BaseEndpoint),
		queue: queue,
	}
}

func (s *HttpClient) Get(path string, resp interface{}) (*http.Response, error) {
	start := time.Now()
	response, err := s.sling.New().Get(path).ReceiveSuccess(&resp)
	elapsed := time.Since(start)

	if response != nil {
		s.queue.Enqueue(entities.NewExternalMetric(response.StatusCode, elapsed.Milliseconds()))
	}

	return response, err
}

func (s *HttpClient) Path(path string) HttpClient {
	s.sling = s.sling.Path(path)

	return *s
}
