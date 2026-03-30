package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const TimeLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	startDate, err := time.Parse(TimeLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart format: %v", err)
	}

	resDate := startDate
	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "y":
		for {
			resDate = resDate.AddDate(1, 0, 0)
			if resDate.After(now) {
				break
			}
		}
	case "d":
		if len(parts) < 2 {
			return "", errors.New("days interval not specified")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days > 400 || days < 1 {
			return "", errors.New("invalid days interval")
		}
		for {
			resDate = resDate.AddDate(0, 0, days)
			if resDate.After(now) {
				break
			}
		}
	default:
		return "", errors.New("unsupported repeat format")
	}

	return resDate.Format(TimeLayout), nil
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(TimeLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}