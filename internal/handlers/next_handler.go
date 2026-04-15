package handlers

import (
	"go-final-project/internal/schedule"
	"net/http"
	"strings"
	"time"
)

func NextHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if nowInput := strings.TrimSpace(r.FormValue("now")); nowInput != "" {
		parsedNow, err := time.Parse(schedule.DateLayout, nowInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsedNow
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	next, err := schedule.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(next))
}
