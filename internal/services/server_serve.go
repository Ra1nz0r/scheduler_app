package services

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	"github.com/ra1nz0r/scheduler_app/internal/models"
)

// Создает адрес запуска сервера и изменяет порт прослушивания по умолчанию,
// если переменная «TODO _PORT» существует в «.env».
// Переменная bool используется один раз, для вывода сообщения о существовании перменной «TODO _PORT» в «.env»
// и изменении стандартного порта при запуске сервера, в остальных случаях пропускается.
func SetServerLink(address string, port string) (string, bool) {
	boolValue := false
	if portFromEnv, exists := os.LookupEnv("TODO_PORT"); exists && portFromEnv != "" {
		port = portFromEnv
		boolValue = true
	}
	return address + port, boolValue
}

// Изменяет путь по умолчанию к базе данных на «TODO_DBFILE», если переменная существует в «.env».
// Переменная bool используется один раз, для вывода сообщения о существовании перменной «TODO_DBFILE» в «.env»
// и изменении стандартного пути датабазы при запуске сервера, в остальных случаях пропускается.
func CheckEnvDbVarOnExists(dbDefaultPath string) (string, bool) {
	boolValue := false
	if dbPathFromEnv := os.Getenv("TODO_DBFILE"); dbPathFromEnv != "" {
		dbDefaultPath, boolValue = dbPathFromEnv, true
	}
	return dbDefaultPath, boolValue
}

// Проверка существования DB.
// Создание папки для хранения DB, файла «.db» и TABLE.
func CheckDBFileExists(resPath string) error {
	if _, errStat := os.Stat(resPath); errStat != nil {
		if os.IsNotExist(errStat) {

			// Создание папки хранения для базы данных.
			folderDb := filepath.Dir(resPath)
			if _, errStat := os.Stat(folderDb); errStat != nil {
				if os.IsNotExist(errStat) {
					if errMkDir := os.Mkdir(folderDb, 0777); errMkDir != nil {
						return fmt.Errorf("failed: cannot create folder: %w", errMkDir)
					}
				}
			}

			logerr.InfoMsg(fmt.Sprintf("Creating %s and TABLE.", filepath.Base(resPath)))
			ctx := context.Background()
			db, errOpen := sql.Open("sqlite", resPath)
			if errOpen != nil {
				logerr.FatalEvent("cannot open DB", errOpen)
			}

			// Создание TABLE.
			if _, errCreate := db.ExecContext(ctx, models.Ddl); errCreate != nil {
				return fmt.Errorf("failed: cannot create table db: %w", errCreate)
			}
			return nil
		}
	}
	logerr.InfoMsg(fmt.Sprintf("Database %s exists.", filepath.Base(resPath)))
	return nil
}
