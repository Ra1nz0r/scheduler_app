package services

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/database/models"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
)

// Проверка существования DB.
// Создание папки для хранения DB, файла «.db» и TABLE.
func CheckDBFileExists(resPath string, db *sql.DB) error {
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
