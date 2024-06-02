package server

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/config"
	"github.com/ra1nz0r/scheduler_app/internal/database"
	hd "github.com/ra1nz0r/scheduler_app/internal/handlers"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	mwar "github.com/ra1nz0r/scheduler_app/internal/middleware"
	"github.com/ra1nz0r/scheduler_app/internal/services"

	"github.com/go-chi/chi"
)

func Run() {
	// Загружаем переменные окружения из '.env' файла.
	conf, errLoad := config.LoadConfig(".")
	if errLoad != nil {
		logerr.FatalEvent("cannot load config", errLoad)
	}

	// Изменяем порт прослушивания по умолчанию, если переменная 'TODO_PORT' существует в '.env'.
	if len(conf.EnvServerPort) > 0 {
		logerr.InfoMsg("'TODO_PORT' exists in '.env' file.")
		logerr.InfoMsg(fmt.Sprintf("Changing default PORT on '%s'.", conf.EnvServerPort))
		config.DefaultPort = conf.EnvServerPort
	}

	// Изменяем путь базы данных по умолчанию, если переменная 'TODO_DBFILE' существует в '.env'»'.
	if len(conf.EnvDatabasePath) > 0 {
		logerr.InfoMsg("'TODO_DBFILE' exists in '.env' file.")
		logerr.InfoMsg(fmt.Sprintf("Changing default PATH on '%s'.", conf.EnvDatabasePath))
		config.DbDefaultPath = conf.EnvDatabasePath
	}

	// ++++++++++++++++
	var d config.DB
	var errOpen error
	d.Db, errOpen = sql.Open("sqlite", config.DbDefaultPath)
	if errOpen != nil {
		logerr.FatalEvent("unable to connect to the database", errOpen)
	}
	var q hd.Queries
	q.Queries = database.New(d.Db)
	// ++++++++++++++++

	logerr.InfoMsg("Checking DB on exists.")
	if errCheck := services.CheckDBFileExists(config.DbDefaultPath, d.Db); errCheck != nil {
		logerr.FatalEvent("cannot check DB on exists", errCheck)
	}

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir(config.DefaultWebDir))
	logerr.InfoMsg("Running handlers.")
	r.Handle("/*", fileServer)

	r.Get("/api/nextdate", hd.NextDateHand)

	r.Get("/api/tasks", mwar.CheckAuth(q.UpcomingTasksWithSearch, conf.EnvPassword))

	r.Post("/api/task/done", mwar.CheckAuth(q.GeneratedNextDate, conf.EnvPassword))

	r.Post("/api/signin", func(w http.ResponseWriter, r *http.Request) {
		hd.LoginAuth(w, r, conf.EnvPassword)
	})

	r.Delete("/api/task", mwar.CheckAuth(q.DeleteTaskScheduler, conf.EnvPassword))
	r.Get("/api/task", mwar.CheckAuth(q.GetTaskByID, conf.EnvPassword))
	r.Post("/api/task", mwar.CheckAuth(q.AddSchedulerTask, conf.EnvPassword))
	r.Put("/api/task", mwar.CheckAuth(q.UpdateTask, conf.EnvPassword))

	serverLink := config.DefIPAddress + ":" + config.DefaultPort

	logerr.InfoMsg(fmt.Sprintf("Starting server on: '%s'", serverLink))

	srv := http.Server{
		Addr:         serverLink,
		Handler:      r,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	go func() {
		if errListn := srv.ListenAndServe(); !errors.Is(errListn, http.ErrServerClosed) {
			logerr.FatalEvent("HTTP server error", errListn)
		}
		logerr.InfoMsg("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if errShut := srv.Shutdown(shutdownCtx); errShut != nil {
		logerr.FatalEvent("HTTP shutdown error", errShut)
	}
	logerr.InfoMsg("Graceful shutdown complete.")
}
