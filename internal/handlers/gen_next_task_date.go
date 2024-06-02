package handlers

import (
	"net/http"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/database"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/services"
)

func (q Queries) GeneratedNextDate(w http.ResponseWriter, r *http.Request) {
	// Получаем задачу по ID и возвращаем ошибку, если её нет в базе данных.
	taskGeted, errGeted := q.Queries.GetTask(r.Context(), r.URL.Query().Get("id"))
	if errGeted != nil {
		ErrReturn(fmt.Errorf("the ID you entered does not exist: %w", errGeted), w)
		return
	}

	switch {
	case taskGeted.Repeat == "": // Одноразовая задача с пустым полем REPEAT удаляется.
		if errDel := q.Queries.DeleteTask(r.Context(), taskGeted.ID); errDel != nil {
			ErrReturn(fmt.Errorf("failed delete: %w", errDel), w)
			return
		}
	default: // В остальных случаях, расчитывается и записывается новая дата для задачи вместо старой.
		newDate, errFunc := services.NextDate(time.Now(), taskGeted.Date, taskGeted.Repeat)
		if errFunc != nil {
			ErrReturn(fmt.Errorf("failed: %w", errFunc), w)
			return
		}

		var task database.UpdateDateTaskParams
		task.ID = taskGeted.ID
		task.Date = newDate
		if errUpd := q.Queries.UpdateDateTask(r.Context(), task); errUpd != nil {
			ErrReturn(fmt.Errorf("failed update task: %w", errUpd), w)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
