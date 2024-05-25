package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Password string `json:"password"`
}

// Обрабатывает POST запрос и принимает на вход пароль пользователя в JSON формате,
// декодирует его в структуру User, в дальнейшем сверяет его с хранящимся в "TODO_PASSWORD"
// в ".env" файле и в случае совпадения отвечатет хэш-суммой пароля. В противном записывает ошибку.
func LoginAuth(w http.ResponseWriter, r *http.Request) {
	// Читаем данные из тела запроса.
	result, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		logerr.ErrEvent("cannot read from BODY", errBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var u User
	if errUnm := json.Unmarshal(result, &u); errUnm != nil {
		services.ErrReturn(fmt.Errorf("can't deserialize: %w", errUnm), w)
		return
	}

	// Проверяем существование переменной "TODO_PASSWORD" в ".env".
	// В случае успеха записываем в результат хэш, в противном ошибку.
	passFromEnv := os.Getenv("TODO_PASSWORD")
	respResult := make(map[string]string)
	switch {
	case passFromEnv == u.Password:
		passHash, errCrypt := bcrypt.GenerateFromPassword([]byte(passFromEnv), bcrypt.DefaultCost)
		if errCrypt != nil {
			services.ErrReturn(fmt.Errorf("failed to generate password hash: %w", errCrypt), w)
		}
		respResult["token"] = string(passHash)
	default:
		services.ErrReturn(fmt.Errorf("incorrect password"), w)
	}

	// Оборачиваем полученные данные в JSON и готовим к выводу,
	// ответ в виде: {"token":"hash"}.
	jsonResp, errJSON := json.Marshal(respResult)
	if errJSON != nil {
		logerr.ErrEvent("failed attempt json-marshal response", errJSON)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusAccepted)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}
