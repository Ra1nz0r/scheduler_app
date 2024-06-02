package services

import (
	"time"

	"fmt"
)

// Условия для 'w'-случаев, задача назначается в указанные дни недели, w <через запятую от 1 до 7> , где 1 — понедельник, 7 — воскресенье.
// Например:
// w 7 — задача перенесётся на ближайшее воскресенье;
// w 1,4,5 — задача перенесётся на ближайший понедельник, четверг или пятницу;
// w 2,3 — задача перенесётся на ближайший вторник или среду.
func weekRepeatCount(clearRep []string, currentDate, startDate time.Time) (string, error) {
	// Получаем числа дней из REPEAT.
	weekDay, errD := RepNumsParse(clearRep[1])
	if errD != nil {
		return "", errD
	}

	// Вычисляем и модифицируем даты в соответствии с переданными в weekDay.
	daysRes, errD := weekModifyDate(weekDay, currentDate, startDate)
	if errD != nil {
		return "", errD
	}

	// Из полученных дат, находим следующую ближайщую после стартовой.
	resDate := FindNearestDate(daysRes, startDate)
	return resDate, nil
}

// Формирует срез следующих дат после стартовой с модифицированными днями, в соответсвии с переданными
// значениями в days. Возвращает ошибку, если число больше или меньше стандартных календарных.
// Формат: [2024-02-29 00:00:00 +0000 UTC 2024-02-18 00:00:00 +0000 UTC]
func weekModifyDate(weekDay []int, currentDate, startDate time.Time) ([]time.Time, error) {
	var daysRes []time.Time
	for _, wNum := range weekDay {
		if wNum < 1 || wNum > 7 {
			return nil, fmt.Errorf("failed: value DAY_WEEK must be between 1 and 7 (your val '%d')", wNum)
		}

		resDate := startDate
		if currentDate.After(startDate) {
			resDate = currentDate
		}

		for ok := true; ok; ok = (resDate.Weekday() != time.Weekday(wNum%7)) { // Находим остаток, потому что Воскресенье = 0, а не 7.
			resDate = resDate.AddDate(0, 0, 1)
		}
		daysRes = append(daysRes, resDate)
	}
	return daysRes, nil
}
