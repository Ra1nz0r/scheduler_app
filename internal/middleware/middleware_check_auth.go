package middleware

import (
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Проверяет существование переменной "TODO_PASSWORD" в ".env" и в положительном случае,
// берёт хэш сумму из cookie, сравнивает с хранящимся паролем и разрешает доступ
// к планировщику. В противном случае, возвращает ошибку и запрещает доступ.
func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passFromEnv, exists := os.LookupEnv("TODO_PASSWORD")
		if exists && passFromEnv != "" {
			var passHash string
			cookie, errCook := r.Cookie("token")
			if errCook == nil {
				passHash = cookie.Value
			}
			if errBC := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(passFromEnv)); errBC != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
