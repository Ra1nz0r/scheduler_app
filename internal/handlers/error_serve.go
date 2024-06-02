package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ra1nz0r/scheduler_app/internal/logerr"
)

// Добавляет ошибки в JSON и возвращает ответ в формате {"error":"ваш текст для ошибки"}.
func ErrReturn(err error, w http.ResponseWriter) {
	result := make(map[string]string)
	result["error"] = err.Error()
	jsonResp, errJSON := json.Marshal(result)
	if errJSON != nil {
		logerr.ErrEvent("failed attempt json-marshal response", errJSON)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusBadRequest)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
