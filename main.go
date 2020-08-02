package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

var Layout = "2006-01-02 15:04:05"

func Now() (r time.Time) {
	location, err := time.LoadLocation("America/Bogota")
	if err != nil {
		errMsg := fmt.Sprintf("Error, LoadLocation in <Now>: %s", "America/Bogota")
		log.Error(errMsg)
	}
	strNow := time.Now().In(location).Format(Layout)
	dateNow, err := ParseDateZ(&strNow)
	if err != nil {
		err := fmt.Sprintf("Error, parsing the date in <shouldNotSkipAssignment>: %s", strNow)
		log.Error(err)
		return
	}

	return *dateNow
}

func ParseDateZ(date *string) (*time.Time, error) {
	if date == nil || *date == "" {
		errMsg := fmt.Sprintf("Error, Empty date in <ParseDateZ>: %d", date)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	location, err := time.LoadLocation("America/Bogota")
	if err != nil {
		log.Errorf("Error, LoadLocation in <ParseDateZ>: %d", date)
		return nil, err
	}

	var res time.Time
	if res, err = time.ParseInLocation(Layout, *date, location); err != nil {
		errMsg := fmt.Sprintf("Error, ParseInLocation in <ParseDateZ>: %d", date)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	return &res, nil
}

func main() {
	//configs := config.NewConfig()
	//redis := redis.NewRedis(configs)
	//val, err := redis.Set("a", "epale", 10000*time.Second)
	//log.Info(val, err)
	//
	//val, err = redis.Get("a")
	//log.Info(val, err)
	//
	//
	//success, err := redis.HSet("b", "epale", "data epoal")
	//log.Info(success, err)
	//
	//val, err = redis.HGet("b", "epale")
	//log.Info(val, err)
	//
	//
	//item := queue.ItemMetric{
	//	Type:         "",
	//	ResponseTime: 0,
	//	StatusCode:   0,
	//	Time:         time.Now(),
	//}
	//
	//log.Info(item)
	//log.Info(item.Field())

}
