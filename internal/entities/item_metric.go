package entities

import (
	"fmt"
	"meli/pkg/slice"
	timeTZ "meli/pkg/time"
	"time"
)

const (
	LocalApi    = "local_api"
	ExternalApi = "external_api"
)

type ItemMetrics struct {
	ResponsesTime slice.Int64
	StatusCode    map[int]int64
	Time          time.Time
}

// ItemMetric the type of the queue
type ItemMetric struct {
	Type         string
	ResponseTime int64
	StatusCode   int
	Time         time.Time
}

func (s *ItemMetric) IsZero() bool {
	return s.Type == ""
}

func (s *ItemMetric) Field() string {
	day, month, year, hour, minute := s.Time.Day(), s.Time.Month().String(), s.Time.Year(), s.Time.Hour(), s.Time.Minute()
	return fmt.Sprintf("%d_%s_%d_%d_%d", day, month, year, hour, minute)
}

func (s *ItemMetric) IsLocal() bool {
	return s.Type == LocalApi
}

func NewLocalMetric(responseTime int64) ItemMetric {
	return ItemMetric{
		Type:         LocalApi,
		ResponseTime: responseTime,
		Time:         timeTZ.Now(),
	}
}

func NewExternalMetric(statusCode int, responseTime int64) ItemMetric {
	return ItemMetric{
		Type:         ExternalApi,
		ResponseTime: responseTime,
		StatusCode:   statusCode,
		Time:         timeTZ.Now(),
	}
}

type ItemMetricResponse struct {
	AvgResponseTime        float64 `json:"avg_response_time"`
	TotalRequests          int64   `json:"total_requests"`
	AvgResponseTimeApiCall float64 `json:"avg_response_time_api_call"`
	TotalCountApiCalls     int64   `json:"total_count_api_calls"`
	Time                   time.Time
	InfoRequests           []InfoRequest `json:"info_requests"`
}

func (s *ItemMetricResponse) IsZero() bool {
	return s.Time.IsZero()
}

type InfoRequest struct {
	StatusCode int64 `json:"status_code"`
	Count      int64 `json:"count"`
}
