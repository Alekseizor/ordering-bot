package conversion

import (
	"strconv"
	"time"
)

const DateFormatLayout = "02.01"

func GetDateStr(today time.Time) string {
	WeekDayIntToStr := map[int]string{0: "Вс", 1: "Пн", 2: "Вт", 3: "Ср", 4: "Чт", 5: "Пт", 6: "Сб"}
	todayStr := strconv.Itoa(today.Day()) + "." + strconv.Itoa(int(today.Month())) + " " + WeekDayIntToStr[int(today.Weekday())]
	return todayStr
}

func GetWeekDayStr(today time.Time) string {
	WeekDayIntToStr := map[int]string{0: "Вс", 1: "Пн", 2: "Вт", 3: "Ср", 4: "Чт", 5: "Пт", 6: "Сб"}
	todayStr := WeekDayIntToStr[int(today.Weekday())]
	return todayStr
}
