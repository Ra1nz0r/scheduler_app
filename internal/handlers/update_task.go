package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/database"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/services"
)

func (q Queries) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Читаем данные из тела запроса.
	result, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		logerr.ErrEvent("cannot read from BODY", errBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var task database.UpdateTaskParams
	if errUnm := json.Unmarshal(result, &task); errUnm != nil {
		ErrReturn(fmt.Errorf("can't deserialize: %w", errUnm), w)
		return
	}

	// Проверяем корректность запроса для обновления параметров задачи в планировщике.
	if _, errFunc := services.NextDate(time.Now(), task.Date, task.Repeat); errFunc != nil {
		ErrReturn(fmt.Errorf("failed: %w", errFunc), w)
		return
	}

	// Если все данные введены корректно, то обновляем задачу в планировщике.
	if errUpdate := q.Queries.UpdateTask(r.Context(), task); errUpdate != nil {
		ErrReturn(fmt.Errorf("can't update task scheduler: %w", errUpdate), w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusAccepted)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
