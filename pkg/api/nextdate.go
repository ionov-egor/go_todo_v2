package api

import (
	"net/http"
	"regexp"
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
