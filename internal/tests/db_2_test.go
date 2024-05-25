package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/ra1nz0r/scheduler_app/internal/config"

	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

func count(db *sqlx.DB) (int, error) {
	var count int
	return count, db.Get(&count, `SELECT count(id) FROM scheduler`)
}

func openDB(t *testing.T) *sqlx.DB {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Error("Error loading .env file")
	}
	dbfile := config.DBFileTest
	envFile := os.Getenv("TODO_DBFILE_TEST")
	if len(envFile) > 0 {
		dbfile = envFile
	}
	db, err := sqlx.Connect("sqlite", dbfile)
	assert.NoError(t, err)
	return db
}
func TestDB(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	before, err := count(db)
	assert.NoError(t, err)

	today := time.Now().Format(`20060102`)

	res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) 
	VALUES (?, 'Todo', 'Комментарий', '')`, today)
	assert.NoError(t, err)

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}

	var task Task
	err = db.Get(&task, `SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?`, id)
	assert.NoError(t, err)
	assert.Equal(t, id, task.ID)
	assert.Equal(t, `Todo`, task.Title)
	assert.Equal(t, `Комментарий`, task.Comment)

	_, err = db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	assert.NoError(t, err)

	after, err := count(db)
	assert.NoError(t, err)

	assert.Equal(t, before, after)
}
