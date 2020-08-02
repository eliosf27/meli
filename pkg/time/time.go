package time

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	Timezone = "America/Bogota"
	Layout   = "2006-01-02 15:04:05"
)

func Now() (r time.Time) {
	location, err := time.LoadLocation(Timezone)
	if err != nil {
		errMsg := fmt.Sprintf("Error, LoadLocation in <Now>: %s", "America/Bogota")
		log.Error(errMsg)

		return time.Time{}
	}

	strNow := time.Now().In(location).Format(Layout)
	var timeResult time.Time
	if timeResult, err = time.ParseInLocation(Layout, strNow, location); err != nil {
		errMsg := fmt.Sprintf("Error, ParseInLocation in <ParseDateZ>: %s", strNow)
		log.Error(errMsg)

		return time.Time{}
	}

	return timeResult
}
