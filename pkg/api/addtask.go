package api

import (
	"strings"
	"time"
)

func getActualDate(date string, repeat string) (string, error) {
	now := time.Now()

	if strings.TrimSpace(date) == "" {
		date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if afterNow(t, now) {
		if len(repeat) == 0 {
			date = now.Format("20060102")
		} else {
			next, err := NextDate(now, date, repeat)
			if err != nil {
				return "", err
			}
			date = next
		}
	}

	return date, nil
}
