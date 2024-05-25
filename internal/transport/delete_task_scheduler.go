package transport

import (
	"context"
	"database/sql"
	"net/http"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/config"
	"github.com/ra1nz0r/scheduler_app/internal/database"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/services"
)

func DeleteTaskScheduler(w http.ResponseWriter, r *http.Request) {
	// Получаем путь из функции и подключаемся к датабазе.
	dbResPath, _ := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	db, errOpen := sql.Open("sqlite", dbResPath)
	if errOpen != nil {
		logerr.FatalEvent("unable to connect to the database", errOpen)
	}

	// Проверям существование задачи и возвращаем ошибку, если её нет в базе данных.
	queries := database.New(db)
	_, errGeted := queries.GetTask(context.Background(), r.URL.Query().Get("id"))
	if errGeted != nil {
		services.ErrReturn(fmt.Errorf("the ID you entered does not exist: %w", errGeted), w)
		return
	}

	// Удаляем задачу из базы данных, при DELETE запросе в виде "/api/task?id=185".
	if errDel := queries.DeleteTask(context.Background(), r.URL.Query().Get("id")); errDel != nil {
		services.ErrReturn(fmt.Errorf("failed delete: %w", errDel), w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
