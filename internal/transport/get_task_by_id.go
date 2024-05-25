package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/config"
	"github.com/ra1nz0r/scheduler_app/internal/database"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/services"
)

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Получаем путь из функции и подключаемся к датабазе.
	dbResPath, _ := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	db, errOpen := sql.Open("sqlite", dbResPath)
	if errOpen != nil {
		logerr.FatalEvent("unable to connect to the database", errOpen)
	}

	// Получаем задачу из планировщика при GET запросе в виде "/api/task?id=185".
	queries := database.New(db)
	idGeted, errGeted := queries.GetTask(context.Background(), r.URL.Query().Get("id"))
	if errGeted != nil {
		services.ErrReturn(fmt.Errorf("the ID you entered does not exist: %w", errGeted), w)
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
