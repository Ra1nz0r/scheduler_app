package transport

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/config"
	"github.com/ra1nz0r/scheduler_app/internal/database"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/services"
)

func GeneratedNextDate(w http.ResponseWriter, r *http.Request) {
	// Получаем путь из функции и подключаемся к базе данных.
	dbResPath, _ := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	db, errOpen := sql.Open("sqlite", dbResPath)
	if errOpen != nil {
		logerr.FatalEvent("unable to connect to the database", errOpen)
	}

	// Получаем задачу по ID и возвращаем ошибку, если её нет в базе данных.
	queries := database.New(db)
	taskGeted, errGeted := queries.GetTask(context.Background(), r.URL.Query().Get("id"))
	if errGeted != nil {
		services.ErrReturn(fmt.Errorf("the ID you entered does not exist: %w", errGeted), w)
		return
	}

	switch {
	case taskGeted.Repeat == "": // Одноразовая задача с пустым полем REPEAT удаляется.
		if errDel := queries.DeleteTask(context.Background(), taskGeted.ID); errDel != nil {
			services.ErrReturn(fmt.Errorf("failed delete: %w", errDel), w)
			return
		}
	default: // В остальных случаях, расчитывается и записывается новая дата для задачи вместо старой.
		newDate, errFunc := services.NextDate(time.Now(), taskGeted.Date, taskGeted.Repeat)
		if errFunc != nil {
			services.ErrReturn(fmt.Errorf("failed: %w", errFunc), w)
			return
		}

		var task database.UpdateDateTaskParams
		task.ID = taskGeted.ID
		task.Date = newDate
		if errUpd := queries.UpdateDateTask(context.Background(), task); errUpd != nil {
			services.ErrReturn(fmt.Errorf("failed update task: %w", errUpd), w)
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
