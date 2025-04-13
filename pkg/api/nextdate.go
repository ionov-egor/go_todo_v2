package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout   = "20060102"
	maxRepeatDay = 400
)

const (
	repeatDay   = "d"
	repeatWeek  = "w"
	repeatMonth = "m"
	repeatYear  = "y"
)

type usedValue struct {
	Key   string
	Value bool
}

var regexpDayInWeek = regexp.MustCompile(`^[1-7]$`)
var regexpDayInMonth = regexp.MustCompile(`^(?:-?[1-2]|[1-9]|0\d|1\d|2\d|3[01])$`)
var regexpMonthNum = regexp.MustCompile(`^[1-9]|1[0-2]$|0[1-9]$`)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	requestNow := r.FormValue(paramNow)
	if requestNow == "" {
		requestNow = time.Now().Format(dateLayout)
	}

	now, err := time.Parse(dateLayout, requestNow)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	newDate, err := NextDate(now, r.FormValue(paramDate), r.FormValue(paramRepeat))
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(newDate)); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func afterNow(now, date time.Time) bool {
	return date.Format("20060102") > now.Format("20060102")
}

func isCorrectNumber(num string, pattren *regexp.Regexp) bool {
	if num == "" {
		return false
	}
	return pattren.MatchString(num)
}

func checkRepeat(repeat string) error {
	if repeat == "" {
		return nil
	}

	subRepeat := strings.Split(repeat, " ")

	switch subRepeat[0] {
	case repeatDay:
		if len(subRepeat) != 2 {
			return errors.New(fmt.Sprintf("Ошибка в поле format: %s. Не корректно заданы значения через пустой символ", repeat))
		}

		interval, err := strconv.Atoi(subRepeat[1])
		if err != nil {
			return err
		}
		if interval > maxRepeatDay {
			return errors.New(fmt.Sprintf("Количество дней повторнеий %d. Превышает максимально допустимое значение %d", interval, maxRepeatDay))
		}

	case repeatWeek:
		if len(subRepeat) != 2 {
			return errors.New(fmt.Sprintf("Ошибка в поле format: %s. Не корректно заданы значения через пустой символ", repeat))
		}

		weekDay := strings.Split(subRepeat[1], ",")
		for key, day := range weekDay {
			if !isCorrectNumber(day, regexpDayInWeek) {
				return errors.New(fmt.Sprintf("Не верно указан день недели %s в позиции %d", day, key))
			}
		}

	case repeatMonth:
		switch len(subRepeat) {
		case 2:
			monthDay := strings.Split(subRepeat[1], ",")
			for key, day := range monthDay {
				if !isCorrectNumber(day, regexpDayInMonth) {
					return errors.New(fmt.Sprintf("Не верно указан день месяца %s в позиции %d", day, key))
				}
			}

		case 3:
			monthDay := strings.Split(subRepeat[1], ",")
			for key, day := range monthDay {
				if !isCorrectNumber(day, regexpDayInMonth) {
					return errors.New(fmt.Sprintf("Не верно указан день месяца %s в позиции %d", day, key))
				}
			}

			monthNum := strings.Split(subRepeat[2], ",")
			for key, day := range monthNum {
				if !isCorrectNumber(day, regexpMonthNum) {
					return errors.New(fmt.Sprintf("Не верно указан месяц %s в позиции %d", day, key))
				}
			}

		default:
			return errors.New(fmt.Sprintf("Ошибка в поле format: %s. Не корректно заданы значения через пустой символ", repeat))
		}

	case repeatYear:
		if len(subRepeat) != 1 {
			return errors.New(fmt.Sprintf("Ошибка в поле format: %s. Не корректно заданы значения через пустой символ", repeat))
		}

	default:
		return errors.New(fmt.Sprintf("Ошибка в поле format: %s. Первый символ может быть только d, w, m или y  ", repeat))
	}

	return nil
}

func getUsedNumc(usedNumc []usedValue, key string) bool {
	if usedNumc == nil {
		return false
	}

	for _, numc := range usedNumc {
		if numc.Key == key {
			return numc.Value
		}
	}

	return false
}

func getNextDateForMonth(daysInMonth []int, daysBeforeEnd []usedValue, months []int, startDate time.Time) (string, error) {
	var date time.Time
	var minDate time.Time
	var isBeforeEnd bool
	var subDate int

	for _, month := range months {
		isBeforeEnd = false
		for _, dayBefore := range daysBeforeEnd {
			if dayBefore.Value {
				if dayBefore.Key == "-1" {
					subDate = -1
				} else {
					subDate = -2
				}

				if isBeforeEnd == false {
					minDate = time.Date(startDate.Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
					minDate = minDate.AddDate(0, 1, 0)
					minDate = minDate.AddDate(0, 0, subDate)
				} else {
					tmpDate := time.Date(startDate.Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
					tmpDate = tmpDate.AddDate(0, 1, 0)
					tmpDate = tmpDate.AddDate(0, 0, subDate)

					if afterNow(tmpDate, minDate) {
						minDate = tmpDate
					}

				}
				isBeforeEnd = true
			}
		}

		if len(daysInMonth) != 0 {
			for _, day := range daysInMonth {
				date = time.Date(startDate.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
				if isBeforeEnd {
					if afterNow(minDate, date) {
						date = minDate
					}
				}

				if date.Day() == day {
					if afterNow(startDate, date) {
						return date.Format("20060102"), nil
					}
				}
			}
		} else {
			if afterNow(startDate, minDate) {
				return minDate.Format("20060102"), nil
			}
		}
	}

	return "", errors.New(fmt.Sprintf("Не удалось выполнить правило повтора"))
}

func nextYear(now time.Time, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(now, date) {
			break
		}
	}

	return date.Format(dateLayout), nil
}

func nextDay(now time.Time, date time.Time, repeat string) (string, error) {
	interval, err := strconv.Atoi(repeat)
	if err != nil {
		return "", err
	}

	for {
		date = date.AddDate(0, 0, interval)
		date.Weekday()

		if afterNow(now, date) {
			break
		}
	}

	return date.Format(dateLayout), nil
}

func nextDayInWeek(now time.Time, date time.Time, repeat string) (string, error) {
	weekDays := make([]usedValue, 0, 7)

	for i := 0; i < 7; i++ {
		weekDays = append(weekDays, usedValue{Key: strconv.Itoa(i + 1)})
	}

	for _, day := range strings.Split(repeat, ",") {
		dayNum, err := strconv.Atoi(day)
		if err != nil {
			return "", err
		}

		weekDays[dayNum-1].Value = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		weekDayInt := int(date.Weekday())
		if weekDayInt == 0 {
			weekDayInt = 7
		}
		if getUsedNumc(weekDays, strconv.Itoa(weekDayInt)) {
			if afterNow(now, date) {
				break
			}
		}
	}

	return date.Format(dateLayout), nil
}

func nextDayInMonth(now time.Time, date time.Time, days string, months string) (string, error) {
	daysInMonth := make([]int, 0, 33)
	daysBeforeEnd := make([]usedValue, 2)

	for _, day := range strings.Split(days, ",") {
		switch day {
		case "-1":
			fallthrough
		case "-2":
			daysBeforeEnd = append(daysBeforeEnd, usedValue{Key: day, Value: true})
		default:
			dayNum, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}
			daysInMonth = append(daysInMonth, dayNum)
		}
	}
	sort.Ints(daysInMonth)

	monthsSlice := make([]int, 0, 12)
	if months == "" {
		for i := 0; i < 12; i++ {
			monthsSlice = append(monthsSlice, i+1)
		}
	} else {
		for _, month := range strings.Split(months, ",") {
			monthNum, err := strconv.Atoi(month)
			if err != nil {
				return "", err
			}
			monthsSlice = append(monthsSlice, monthNum)
		}
	}
	sort.Ints(monthsSlice)

	if afterNow(date, now) {
		date = now
	}

	newDate, err := getNextDateForMonth(daysInMonth, daysBeforeEnd, monthsSlice, date)
	if err != nil {
		return "", err
	}

	return newDate, nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", nil
	}

	if err := checkRepeat(repeat); err != nil {
		return "", err
	}

	subRepeat := strings.Split(repeat, " ")

	switch subRepeat[0] {
	case repeatYear:
		newDate, err := nextYear(now, date)
		if err != nil {
			return "", err
		}
		return newDate, nil

	case repeatDay:
		newDate, err := nextDay(now, date, subRepeat[1])
		if err != nil {
			return "", err
		}
		return newDate, nil
	case repeatWeek:
		newDate, err := nextDayInWeek(now, date, subRepeat[1])
		if err != nil {
			return "", err
		}
		return newDate, nil

	case repeatMonth:
		var monthsNum string
		if len(subRepeat) == 3 {
			monthsNum = subRepeat[2]
		} else {
			monthsNum = ""
		}

		newDate, err := nextDayInMonth(now, date, subRepeat[1], monthsNum)
		if err != nil {
			return "", err
		}
		return newDate, nil
	}

	return "", nil
}
