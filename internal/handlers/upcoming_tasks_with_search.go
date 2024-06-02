package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/logerr"
)

// Обработчик для GET-запроса /api/tasks. Возвращает список ближайших задач в
// формате JSON в виде списка в поле tasks.
// Также обрабатывает параметр search в строке запроса "/api/tasks?search=бассейн",
// в данном случае возвратит задачи со словом «бассейн».
// Включая поиск задач по указанной дате в формате "/api/tasks?search=08.02.2024".
func (q Queries) UpcomingTasksWithSearch(w http.ResponseWriter, r *http.Request) {
	// Создание мапы для выведения полученных данных в виде:
	// {"tasks":[{задача1}, {задача2}, {задача3}...}]}.
	// Используется интерфейс, потому что записываются разные структуры.
	respResult := make(map[string]interface{})

	// Получение списка ближайших задач.
	switch {
	// В случае, если запрашивается поиск "/api/tasks?search={данные}".
	case len(strings.TrimSpace(r.URL.Query().Get("search"))) > 0:
		// Обрабатываем и проверяем дату на соответствие формату.
		dateSearch, errSearch := time.Parse("02.01.2006", r.URL.Query().Get("search"))
		switch {
		// Если дата в запросе соответствует формату, то выполняем SEARCH в датабазе по запрашиваемой дате.
		case errSearch == nil:
			resDate, resErr := q.Queries.SearchDate(r.Context(), dateSearch.Format("20060102"))
			if resErr != nil {
				ErrReturn(fmt.Errorf("%w", resErr), w)
				return
			}
			respResult["tasks"] = resDate
			// Чтобы избежать {"tasks":null} в ответе JSON при отсутствии результат,
			// перезаписываем полученные данные на {"tasks": []}.
			if resDate == nil {
				respResult["tasks"] = []string{}
			}
		// Если дата в запросе НЕ СООТВЕТСТВУЕТ формату, то выполняем SEARCH в датабазе по запрошенному слову.
		default:
			// Приводим к нижнему регистру и выполняем SEARCH по такой же колонке в датабазе.
			search := strings.ToLower("%" + r.URL.Query().Get("search") + "%")
			resSearch, resErr := q.Queries.SearchTasks(r.Context(), search)
			if resErr != nil {
				ErrReturn(fmt.Errorf("%w", resErr), w)
				return
			}
			respResult["tasks"] = resSearch
			// Чтобы избежать {"tasks":null} в ответе JSON при отсутствии результат,
			// перезаписываем полученные данные на {"tasks": []}.
			if resSearch == nil {
				respResult["tasks"] = []string{}
			}
		}
	// Если запрос получен в стандартной форме "/api/tasks", то возвращаем
	// список ближайших задач в формате JSON в виде списка в поле tasks.
	default:
		resList, resErr := q.Queries.ListTasks(r.Context())
		if resErr != nil {
			ErrReturn(fmt.Errorf("%w", resErr), w)
			return
		}
		respResult["tasks"] = resList
		// Чтобы избежать {"tasks":null} в ответе JSON при отсутствии результат,
		// перезаписываем полученные данные на {"tasks": []}.
		if resList == nil {
			respResult["tasks"] = []string{}
		}
	}

	// Оборачиваем полученные данные в JSON и готовим к выводу,
	// ответ в виде: {"tasks":[{task1}, {task2}, .... ]}.
	jsonResp, errJSON := json.Marshal(respResult)
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
