package conversion

import (
	"strconv"
	"time"
)

func GetDateStr(today time.Time) string {
	var dayStr, monStr string
	WeekDayIntToStr := map[int]string{0: "Вс", 1: "Пн", 2: "Вт", 3: "Ср", 4: "Чт", 5: "Пт", 6: "Сб"}
	if today.Day() < 10 {
		dayStr = "0" + strconv.Itoa(today.Day())
	} else {
		dayStr = strconv.Itoa(today.Day())
	}
	if int(today.Month()) < 10 {
		monStr = "0" + strconv.Itoa(int(today.Month()))
	} else {
		monStr = strconv.Itoa(int(today.Month()))
	}
	todayStr := dayStr + "." + monStr + " " + WeekDayIntToStr[int(today.Weekday())]
	return todayStr
}

func GetWeekDayStr(today time.Time) string {
	WeekDayIntToStr := map[int]string{0: "Вс", 1: "Пн", 2: "Вт", 3: "Ср", 4: "Чт", 5: "Пт", 6: "Сб"}
	todayStr := WeekDayIntToStr[int(today.Weekday())]
	return todayStr
}
