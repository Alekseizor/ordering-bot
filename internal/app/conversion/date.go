package conversion

import (
	"time"
)

const DateFormatLayout = "02.01"

func GetUpcomingDates(count int) []string {
	upcomingDates := make([]string, count, count)
	t := time.Now()

	t = t.Add(time.Hour * 48)

	for i := 0; i < count; i++ {
		upcomingDates[i] = t.Format(DateFormatLayout)
		t = t.Add(time.Hour * 24)
	}

	return upcomingDates
}
