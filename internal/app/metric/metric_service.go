package metric

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"meli/internal/entities"
	"meli/internal/redis"
	"meli/pkg/slice"
)

const MetricsKey = "meli:metrics"

type MetricServicer interface {
	UpdateMetric(item entities.ItemMetric) error
	FetchMetrics() []entities.ItemMetricResponse
}

type MetricService struct {
	redis redis.Redis
}

func NewMetricService(redis redis.Redis) MetricServicer {
	return &MetricService{redis: redis}
}

func (s *MetricService) UpdateMetric(item entities.ItemMetric) error {
	key := MetricsKey
	field := item.Field()
	metrics, _ := s.fetch(key, field)
	metrics = s.calculate(item, metrics)

	return s.save(key, field, metrics)
}

func (s *MetricService) FetchMetrics() []entities.ItemMetricResponse {
	metricsMap, err := s.redis.HGetAll(MetricsKey)
	if err != nil {
		log.Errorf("error getting all metrics [%s] - %+v", MetricsKey, err)

		return []entities.ItemMetricResponse{}
	}

	return s.buildMetrics(metricsMap)
}

func (s *MetricService) buildMetrics(metricsMap map[string]string) []entities.ItemMetricResponse {
	var responses []entities.ItemMetricResponse

	for _, val := range metricsMap {
		metricResponse := entities.ItemMetricResponse{
			InfoRequests: []entities.InfoRequest{},
		}
		currentMetrics := map[string]entities.ItemMetrics{}
		if err := json.Unmarshal([]byte(val), &currentMetrics); err != nil {
			log.Errorf("error unmarshal currentMetrics: %+v", err)

			return responses
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
			metricResponse.InfoRequests = s.buildInfo(externalMetric.StatusCode)
		}

		if !metricResponse.IsZero() {
			responses = append(responses, metricResponse)
		}
	}

	return responses
}

func (s *MetricService) buildInfo(statusCodes map[int]int64) []entities.InfoRequest {
	var infos []entities.InfoRequest

	for statusCode, count := range statusCodes {
		info := entities.InfoRequest{
			StatusCode: int64(statusCode),
			Count:      count,
		}

		infos = append(infos, info)
	}

	return infos
}

func (s *MetricService) fetch(key string, field string) (map[string]entities.ItemMetrics, error) {
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

func (s *MetricService) save(key string, field string, metrics map[string]entities.ItemMetrics) error {
	metricsRaw, err := json.Marshal(&metrics)
	if err != nil {
		log.Errorf(
			"Error marshalling item for saving metrics [%s] - %+v", key, err,
		)

		return err
	}

	_, err = s.redis.HSet(key, field, metricsRaw)
	if err != nil {
		log.Errorf("Error saving metrics [%s] - %+v", key, err)

		return err
	}

	return nil
}

func (s *MetricService) calculate(item entities.ItemMetric, metrics map[string]entities.ItemMetrics) map[string]entities.ItemMetrics {
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
