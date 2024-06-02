package handlers

import (
	"net/http"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/logerr"
)

func (q Queries) DeleteTaskScheduler(w http.ResponseWriter, r *http.Request) {
	// Проверям существование задачи и возвращаем ошибку, если её нет в базе данных.
	_, errGeted := q.Queries.GetTask(r.Context(), r.URL.Query().Get("id"))
	if errGeted != nil {
		ErrReturn(fmt.Errorf("the ID you entered does not exist: %w", errGeted), w)
		return
	}

	// Удаляем задачу из базы данных, при DELETE запросе в виде "/api/task?id=185".
	if errDel := q.Queries.DeleteTask(r.Context(), r.URL.Query().Get("id")); errDel != nil {
		ErrReturn(fmt.Errorf("failed delete: %w", errDel), w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
