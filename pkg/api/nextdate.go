package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout   = "20060102"
	repeatFormat = "dywm"
	repeatDay    = "d"
	maxRepeatDay = 400
)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue(paramNow)
	if now == "" {
		now = time.Now().Format(dateLayout)
	}

	nowTime, err := time.Parse(dateLayout, now)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	newDate, err := NextDate(nowTime, r.FormValue(paramDate), r.FormValue(paramRepeat))
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newDate))
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func checkRepeat(repeat string) error {
	subRepeat := strings.Split(repeat, " ")

	switch len(subRepeat) {
	case 1:
		if subRepeat[0] != "y" {
			return errors.New(fmt.Sprintf("Ошибка в поле format: %s", repeat))
		}
	case 2:
	default:
		return errors.New(fmt.Sprintf("Ошибка в поле format: %s", repeat))
	}

	if strings.Contains(repeatFormat, subRepeat[0]) == false {
		return errors.New(fmt.Sprintf("Ошибка в поле format: %s", repeat))
	}

	if subRepeat[0] == repeatDay {
		interval, err := strconv.Atoi(subRepeat[1])
		if err != nil {
			return err
		}
		if interval > maxRepeatDay {
			return errors.New(fmt.Sprintf("Количество дней повторнеий %d. Превышает максимально допустимое значение %d", interval, maxRepeatDay))
		}

	}

	return nil
}

func nextYear(now time.Time, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateLayout), nil
}

func nextDay(now time.Time, date time.Time, repeat string) (string, error) {
	if now.Format(dateLayout) == date.Format(dateLayout) {
		return date.Format(dateLayout), nil
	}

	interval, err := strconv.Atoi(repeat)
	if err != nil {
		return "", err
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateLayout), nil
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
	case "y":
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

	case "m":
		fallthrough
	case "w":
		return "", errors.New(fmt.Sprintf("Ошибка в поле format: %s", subRepeat[0]))
	}

	return "", nil
}
