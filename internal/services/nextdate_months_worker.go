package services

import (
	"time"

	"fmt"
)

// Для 'm'-случаев, задача назначается в указанные дни недели, m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12].
// При этом вторая последовательность чисел опциональна и указывает на определённые месяцы.
// Например:
// m 4 — задача назначается на четвёртое число каждого месяца;
// m 1,15,25 — задача назначается на 1-е, 15-е и 25-е число каждого месяца;
// m -1 — задача назначается на последний день месяца;
// m -2 — задача назначается на предпоследний день месяца;
// m 3 1,3,6 — задача назначается на 3-е число января, марта и июня;
// m 1,-1 2,8 — задача назначается на 1-е и последнее число число февраля и авгуcта.
func mounthRepeatCount(clearRep []string, currentDate, startDate time.Time) (string, error) {
	var resDate string
	var errRes error
	switch {
	case len(clearRep) == 2:
		resDate, errRes = mRepeatWithout(clearRep, currentDate, startDate)
		if errRes != nil {
			return "", errRes
		}
	// Вычисления для m-случаев, с переданными днями месяцев и с указанием конкретных месяцев.
	case len(clearRep) == 3:
		resDate, errRes = mRepeatWithMonths(clearRep, currentDate, startDate)
		if errRes != nil {
			return "", errRes
		}
	}
	return resDate, nil
}

// === Для случаев без месяцев в REPEAT: "m 5", "m 10,17", ... === //

func mRepeatWithout(clearRep []string, currentDate, startDate time.Time) (string, error) {
	// Получаем числа дней из REPEAT.
	monthDays, errD := RepNumsParse(clearRep[1])
	if errD != nil {
		return "", errD
	}

	// Вычисляем и модифицируем даты в соответствии с переданными в monthDays.
	modiDateRes, errD := modifyDate(monthDays, currentDate, startDate)
	if errD != nil {
		return "", errD
	}

	// Если текущая дата идет после стартовой, меняем значение расчётной на текущую.
	if currentDate.After(startDate) {
		startDate = currentDate
	}

	// Из полученных дат, находим следующую ближайщую после стартовой.
	resDate := FindNearestDate(modiDateRes, startDate)
	return resDate, nil
}

// Формирует срез следующих дат после стартовой с модифицированными днями, в соответсвии с переданными
// значениями в days. Возвращает ошибку, если число больше или меньше стандартных календарных.
// При передаче -1 и -2 вычисляется последний и предпоследний день месяца.
// Формат: [2024-02-29 00:00:00 +0000 UTC 2024-02-18 00:00:00 +0000 UTC]
func modifyDate(days []int, currentDate, startDate time.Time) ([]time.Time, error) {
	var daysRes []time.Time
	for _, dNum := range days {
		if dNum < -2 || dNum > 31 {
			return nil, fmt.Errorf("failed: value DAY_MONTH must be between -2 and 31 (your val '%d')", dNum)
		}

		resDate := startDate
		if currentDate.After(startDate) {
			resDate = currentDate
		}

		var tt time.Time
		switch dNum {
		case -1:
			// В случае, если дата равна последнему дню месяца, то чтобы не получить такую же дату, а следующую.
			// Мы прибавляем один день к дате, затем прибавив один месяц получем последний день следующего месяца.
			resDate = resDate.AddDate(0, 0, 1)
			tt = time.Date(resDate.Year(), resDate.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			daysRes = append(daysRes, tt)
		case -2:
			// Аналогично случаю "-1", только для предпоследнего дня. Избегаем получения такой же даты.
			resDate = resDate.AddDate(0, 0, 2)
			tt = time.Date(resDate.Year(), resDate.Month()+1, -1, 0, 0, 0, 0, time.UTC)
			daysRes = append(daysRes, tt)
		default:
			for ok := true; ok; ok = (resDate.Day() != dNum) {
				resDate = resDate.AddDate(0, 0, 1)
			}
			daysRes = append(daysRes, resDate)
		}
	}
	return daysRes, nil
}

// === Для случаев с месяцами в REPEAT: "m 5 1,13", "m 10,17 12,8,1", ... === //

func mRepeatWithMonths(clearRep []string, currentDate, startDate time.Time) (string, error) {
	// Получаем месяца из REPEAT.
	months, errM := RepNumsParse(clearRep[2])
	if errM != nil {
		return "", errM
	}
	// Вычисляем даты месяцев из переданных в REPAT и меняем день на первый.
	monthsDateRes, errM := findingMonth(months, currentDate, startDate)
	if errM != nil {
		return "", errM
	}
	// Получаем числа дней из REPEAT.
	monthDays, errD := RepNumsParse(clearRep[1])
	if errD != nil {
		return "", errD
	}
	// Модифицируем даты monthRes, изменяя дни на переданные в monthDays.
	modiDateRes, errD := modifyDayMonth(monthsDateRes, monthDays)
	if errD != nil {
		return "", errD
	}

	// Если текущая дата идет после стартовой, меняем значение расчётной на текущую.
	if currentDate.After(startDate) {
		startDate = currentDate
	}

	// Из полученных дат, находим следующую ближайщую после стартовой.
	resDate := FindNearestDate(modiDateRes, startDate)
	return resDate, nil
}

// Принимает срез чисел месяцев [12, 8, ... ], вычисляет даты соответствующие этим месяцам, изменяя
// день на первый. И добавляет эти месяца в возращаемый срез дат monthsRes.
// Возвращает ошибку, если число больше или меньше стандартных календарных.
// Формат: [2024-12-01 00:00:00 +0000 UTC 2024-08-01 00:00:00 +0000 UTC ... ]
func findingMonth(months []int, currentDate, startDate time.Time) ([]time.Time, error) {
	var monthsRes []time.Time
	for _, mNum := range months {
		if mNum < 1 || mNum > 12 {
			return nil, fmt.Errorf("failed: value MONTH must be between 1 and 12 (your val '%d')", mNum)
		}

		resDate := startDate
		if currentDate.After(resDate) {
			resDate = currentDate
		}

		for {
			if resDate.Month() == time.Month(mNum) {
				break
			}
			resDate = resDate.AddDate(0, 1, 0)
		}
		tt := time.Date(resDate.Year(), resDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		monthsRes = append(monthsRes, tt)
	}
	return monthsRes, nil
}

// Модифицирует из принимаемого среза monthsRes все даты, изменяя день на передаваемые числа в срезе days.
// Возвращает ошибку, если число больше или меньше стандартных календарных.
// При передаче -1 и -2 вычисляется последний и предпоследний день месяца.
// Возвращает модифицированный срез дат из переданных в monthRes.
// Формат: [2024-12-10 00:00:00 +0000 UTC 2024-12-17 00:00:00 +0000 ... ]
func modifyDayMonth(monthsRes []time.Time, days []int) ([]time.Time, error) {
	var daysRes []time.Time
	for _, mNum := range monthsRes {
		for _, dNum := range days {
			if dNum < -2 || dNum > 31 {
				return nil, fmt.Errorf("failed: value DAY_MONTH must be between -2 and 31 (your val '%d')", dNum)
			}
			var tt time.Time
			switch dNum {
			case -1:
				tt = time.Date(mNum.Year(), mNum.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			case -2:
				tt = time.Date(mNum.Year(), mNum.Month()+1, -1, 0, 0, 0, 0, time.UTC)
			default:
				tt = time.Date(mNum.Year(), mNum.Month(), dNum, 0, 0, 0, 0, time.UTC)
			}
			daysRes = append(daysRes, tt)
		}
	}
	return daysRes, nil
}
