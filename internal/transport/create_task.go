package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/config"
	"github.com/ra1nz0r/scheduler_app/internal/database"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/services"
)

// Обработчик для POST запросов и добавления задачи в датабазу. Запрос и ответ передаются в JSON-формате.
// Запрос состоит из следующих "string" полей:
// date — дата задачи в формате 20060102;
// title — заголовок задачи, обязательное поле;
// comment — комментарий к задаче;
// repeat — правило повторения задачи.
func AddSchedulerTask(w http.ResponseWriter, r *http.Request) {
	// Читаем данные из тела запроса.
	result, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		logerr.ErrEvent("cannot read from BODY", errBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var task database.CreateTaskParams
	if errUnm := json.Unmarshal(result, &task); errUnm != nil {
		services.ErrReturn(fmt.Errorf("can't deserialize: %w", errUnm), w)
		return
	}

	// Проверка на отсутствие поля TITLE.
	if len(strings.TrimSpace(task.Title)) == 0 {
		services.ErrReturn(fmt.Errorf("failed: TITLE cannot be EMPTY"), w)
		return
	}

	// Если DATE не заполнена при добавлении задачи, то перезаписываем на текущую.
	if len(strings.TrimSpace(task.Date)) == 0 {
		task.Date = time.Now().Format("20060102")
	}

	// Проверка корректности даты.
	if _, errPars := time.Parse("20060102", task.Date); errPars != nil {
		services.ErrReturn(fmt.Errorf("failed, incorrect DATE: %w", errPars), w)
		return
	}

	// Если введеная DATE меньше текущей и поле REPEAT не заполнено,
	// перезаписываем DATE на текущую дату. В противном случае, проверяем корректность
	// REPEAT и если DATE меньше текущей, то перезаписываем на res значение.
	if task.Date < time.Now().Format("20060102") {
		switch {
		case len(strings.TrimSpace(task.Repeat)) == 0:
			task.Date = time.Now().Format("20060102")
		default:
			res, errFunc := services.NextDate(time.Now(), task.Date, task.Repeat)
			if errFunc != nil {
				services.ErrReturn(fmt.Errorf("failed: %w", errFunc), w)
				return
			}
			if task.Date < time.Now().Format("20060102") {
				task.Date = res
			}
		}
	}

	// Получаем путь из функции и подключаемся к датабазе.
	dbResPath, _ := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	db, errOpen := sql.Open("sqlite", dbResPath)
	if errOpen != nil {
		logerr.FatalEvent("unable to connect to the database", errOpen)
	}

	// Если данные корректны, то создаём запись в датабазе.
	queries := database.New(db)
	insertedTask, errCreate := queries.CreateTask(context.Background(), task)
	if errCreate != nil {
		logerr.ErrEvent("cannot create task in DB", errCreate)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Создание мапы и выведение последнего ID добавленного в датабазу, ответ в виде: {"id":"186"}.
	respResult := make(map[string]string)
	respResult["id"] = insertedTask.ID
	jsonResp, errJSON := json.Marshal(respResult)
	if errJSON != nil {
		logerr.ErrEvent("failed attempt json-marshal response", errJSON)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusCreated)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
