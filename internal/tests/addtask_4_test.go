package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"testing"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/config"
	"github.com/stretchr/testify/assert"
)

func requestJSON(apipath string, values map[string]any, method string) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if len(values) > 0 {
		data, err = json.Marshal(values)
		if err != nil {
			return nil, err
		}
	}
	var resp *http.Response

	req, err := http.NewRequest(method, getURL(apipath), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	// Загружаем переменные окружения из '.env' файла.
	conf, errLoad := config.LoadConfig("../..")
	if errLoad != nil {
		log.Fatal("No .env file found")
	}
	if len(conf.EnvPassHashForTest) > 0 {
		jar, errJar := cookiejar.New(nil)
		if errJar != nil {
			return nil, errJar
		}
		jar.SetCookies(req.URL, []*http.Cookie{
			{
				Name:  "token",
				Value: conf.EnvPassHashForTest,
			},
		})
		client.Jar = jar
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return io.ReadAll(resp.Body)
}

func postJSON(apipath string, values map[string]any, method string) (map[string]any, error) {
	var (
		m   map[string]any
		err error
	)

	body, err := requestJSON(apipath, values, method)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &m)
	return m, err
}

type task struct {
	date    string
	title   string
	comment string
	repeat  string
}

func TestAddTask(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	tbl := []task{
		{"20240129", "", "", ""},
		{"20240192", "Qwerty", "", ""},
		{"28.01.2024", "Заголовок", "", ""},
		{"20240112", "Заголовок", "", "w"},
		{"20240212", "Заголовок", "", "ooops"},
	}
	for _, v := range tbl {
		m, err := postJSON("api/task", map[string]any{
			"date":    v.date,
			"title":   v.title,
			"comment": v.comment,
			"repeat":  v.repeat,
		}, http.MethodPost)
		assert.NoError(t, err)

		e, ok := m["error"]
		assert.False(t, !ok || len(fmt.Sprint(e)) == 0,
			"Ожидается ошибка для задачи %v", v)
	}

	now := time.Now()

	check := func() {
		for _, v := range tbl {
			today := v.date == "today"
			if today {
				v.date = now.Format(`20060102`)
			}
			m, err := postJSON("api/task", map[string]any{
				"date":    v.date,
				"title":   v.title,
				"comment": v.comment,
				"repeat":  v.repeat,
			}, http.MethodPost)
			assert.NoError(t, err)

			e, ok := m["error"]
			if ok && len(fmt.Sprint(e)) > 0 {
				t.Errorf("Неожиданная ошибка %v для задачи %v", e, v)
				continue
			}
			var task Task
			var mid any
			mid, ok = m["id"]
			if !ok {
				t.Errorf("Не возвращён id для задачи %v", v)
				continue
			}
			id := fmt.Sprint(mid)

			err = db.Get(&task, `SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?`, id)
			assert.NoError(t, err)
			assert.Equal(t, id, strconv.FormatInt(task.ID, 10))

			assert.Equal(t, v.title, task.Title)
			assert.Equal(t, v.comment, task.Comment)
			assert.Equal(t, v.repeat, task.Repeat)
			if task.Date < now.Format(`20060102`) {
				t.Errorf("Дата не может быть меньше сегодняшней %v", v)
				continue
			}
			if today && task.Date != now.Format(`20060102`) {
				t.Errorf("Дата должна быть сегодняшняя %v", v)
			}
		}
	}

	tbl = []task{
		{"", "Заголовок", "", ""},
		{"20231220", "Сделать что-нибудь", "Хорошо отдохнуть", ""},
		{"20240108", "Уроки", "", "d 10"},
		{"20240102", "Отдых в Сочи", "На лыжах", "y"},
		{"today", "Фитнес", "", "d 1"},
		{"today", "Шмитнес", "", ""},
	}
	check()
	if config.FullNextDate {
		tbl = []task{
			{"20240129", "Сходить в магазин", "", "w 1,3,5"},
		}
		check()
	}
}
