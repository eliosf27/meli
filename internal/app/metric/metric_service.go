package metric

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"meli/internal/app/entities"
	"meli/internal/redis"
	"meli/pkg/slice"
)

const MetricsKey = "meli-metrics"

type MetricService interface {
	UpdateMetric(item entities.ItemMetric) error
	FetchMetrics() []entities.ItemMetricResponse
}

type service struct {
	redis redis.Redis
}

func NewMetricService(redis redis.Redis) MetricService {
	return &service{redis: redis}
}

func (s *service) UpdateMetric(item entities.ItemMetric) error {
	key := MetricsKey
	field := item.Field()
	metrics, err := s.fetch(key, field)
	if err != nil {
		return err
	}

	metrics = s.calculate(item, metrics)

	return s.save(key, field, metrics)
}

func (s *service) FetchMetrics() []entities.ItemMetricResponse {
	var responses []entities.ItemMetricResponse
	metricsMap, err := s.redis.HGetAll(MetricsKey)
	if err != nil {
		log.Errorf("error getting all metrics [%s] - %+v", MetricsKey, err)

		return responses
	}

	for _, val := range metricsMap {
		metricResponse := entities.ItemMetricResponse{
			InfoRequests: []entities.InfoRequest{},
		}
		currentMetrics := map[string]entities.ItemMetrics{}
		if err := json.Unmarshal([]byte(val), &currentMetrics); err != nil {
			log.Errorf("Error unmarshal currentMetrics: %+v", err)
		}

		if localMetric, ok := currentMetrics[entities.LocalApi]; ok {
			metricResponse.Time = localMetric.Time
			metricResponse.TotalRequests = localMetric.ResponsesTime.Count()
			metricResponse.AvgResponseTime = localMetric.ResponsesTime.Avg()
		}

		if externalMetric, ok := currentMetrics[entities.ExternalApi]; ok {
			metricResponse.Time = externalMetric.Time
			metricResponse.TotalCountApiCalls = externalMetric.ResponsesTime.Count()
			metricResponse.AvgResponseTimeApiCall = externalMetric.ResponsesTime.Avg()

			var infos []entities.InfoRequest
			for statusCode, count := range externalMetric.StatusCode {
				info := entities.InfoRequest{
					StatusCode: int64(statusCode),
					Count:      count,
				}

				infos = append(infos, info)
			}

			metricResponse.InfoRequests = infos
		}

		if !metricResponse.IsZero() {
			responses = append(responses, metricResponse)
		}
	}

	return responses
}

func (s *service) fetch(key string, field string) (map[string]entities.ItemMetrics, error) {
	metrics := map[string]entities.ItemMetrics{}
	val, err := s.redis.HGet(key, field)
	if err != nil {
		log.Errorf("Error fetching metric [%s] - %+v", key, err)

		return metrics, err
	}

	if err := json.Unmarshal(val, &metrics); err != nil {
		log.Errorf("Error unmarshal: %+v", err)

		return metrics, err
	}

	return metrics, nil
}

func (s *service) save(key string, field string, metrics map[string]entities.ItemMetrics) error {
	metricsRaw, err := json.Marshal(&metrics)
	if err != nil {
		log.Errorf(
			"Error marshalling item for saving metrics [%s] - [%s] - %+v", key, err,
		)

		return err
	}

	_, err = s.redis.HSet(key, field, metricsRaw)
	if err != nil {
		log.Errorf("Error saving metrics [%s] - [%s] - %+v", key, err)

		return err
	}

	return nil
}

func (s *service) calculate(item entities.ItemMetric, metrics map[string]entities.ItemMetrics) map[string]entities.ItemMetrics {
	metric := entities.ItemMetrics{
		ResponsesTime: slice.SliceInt64{},
		StatusCode:    map[int]int64{},
		Time:          item.Time,
	}
	if currentMetric, ok := metrics[item.Type]; ok {
		metric = currentMetric
	}

	metric.ResponsesTime = append(metric.ResponsesTime, item.ResponseTime)
	if !item.IsLocal() {
		metric.StatusCode[item.StatusCode] += 1
	}

	metrics[item.Type] = metric

	return metrics
}