package services

import (
	"time"
)

// Вычисляет следующую дату для запроса REPEAT с буквой 'y' - задача переносится год.
func yearRepeatCount(currentDate, startDate time.Time) string {
	switch startDate.After(currentDate) {
	case true:
		startDate = startDate.AddDate(1, 0, 0)
	default:
		for currentDate.After(startDate) {
			startDate = startDate.AddDate(1, 0, 0)
		}
	}
	return startDate.Format("20060102")
}
