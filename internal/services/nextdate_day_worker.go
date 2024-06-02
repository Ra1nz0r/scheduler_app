package services

import (
	"time"

	"fmt"
)

// Условия для 'd'-случаев, задача переносится на указанное число дней, максимально допустимое число равно 400.
// Например :
// d 1 — каждый день;
// d 7 — для вычисления следующей даты добавляем семь дней;
// d 60 — переносим на 60 дней.
func dayRepeatCount(clearRep []string, currentDate, startDate time.Time) (string, error) {
	// Получаем числа дней из REPEAT.
	days, errD := RepNumsParse(clearRep[1])
	if errD != nil {
		return "", errD
	}

	// Вычисляем и модифицируем даты в соответствии с переданными в days.
	resDate, errD := dayResWork(days, currentDate, startDate)
	if errD != nil {
		return "", errD
	}
	return resDate, nil
}

func dayResWork(days []int, currentDate, startDate time.Time) (string, error) {
	if days[0] < 1 || days[0] > 400 {
		return "", fmt.Errorf("failed: value DAY must be between 1 and 400 (your val '%d')", days[0])
	}
	switch startDate.After(currentDate) {
	case true:
		startDate = startDate.AddDate(0, 0, days[0])
	default:
		for currentDate.After(startDate) {
			startDate = startDate.AddDate(0, 0, days[0])
		}
	}
	return startDate.Format("20060102"), nil
}
