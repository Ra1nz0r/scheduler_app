package handlers

import (
	"net/http"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/services"
)

// Обработчик для GET запросов и вывода следующей даты в текстовом формате
// для задачи в планировщике, после обработки функциeй NextDate.
func NextDateHand(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем введенную дату.
	todayDate, errPars := time.Parse("20060102", r.URL.Query().Get("now"))
	if errPars != nil {
		errorMsg := fmt.Sprintf("Failed: incorrect DATE: %v", errPars)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	// Вычисление следующей даты, подробнее в описании NextDate.
	res, errFunc := services.NextDate(todayDate, r.URL.Query().Get("date"), r.URL.Query().Get("repeat"))
	if errFunc != nil {
		http.Error(w, errFunc.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(res)); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
