package handlers

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/logerr"
)

func (q Queries) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Получаем задачу из планировщика при GET запросе в виде "/api/task?id=185".
	idGeted, errGeted := q.Queries.GetTask(r.Context(), r.URL.Query().Get("id"))
	if errGeted != nil {
		ErrReturn(fmt.Errorf("the ID you entered does not exist: %w", errGeted), w)
		return
	}

	// Оборачиваем полученные данные в JSON и готовим к выводу,
	// ответ в виде: {"id": "айди","date": "дата","title": "заголовок","comment": "коммент","repeat": "условия повторения"}.
	jsonResp, errJSON := json.Marshal(idGeted)
	if errJSON != nil {
		logerr.ErrEvent("failed attempt json-marshal response", errJSON)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
