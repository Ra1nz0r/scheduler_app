package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"github.com/ra1nz0r/scheduler_app/internal/config"
	"github.com/ra1nz0r/scheduler_app/internal/logerr"
	mwar "github.com/ra1nz0r/scheduler_app/internal/middleware"
	"github.com/ra1nz0r/scheduler_app/internal/services"
	tp "github.com/ra1nz0r/scheduler_app/internal/transport"

	"github.com/go-chi/chi"
)

func Run() {
	serverLink, boolValue := services.SetServerLink("0.0.0.0:", config.DefaultPort)
	if boolValue {
		logerr.InfoMsg("'TODO_PORT' exists in '.env' file. Changing default PORT.")
	}
	dbResultPath, boolValue := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	if boolValue {
		logerr.InfoMsg("'TODO_DBFILE' exists in '.env' file. Changing default PATH.")
	}

	logerr.InfoMsg("Checking DB on exists.")
	if errCheck := services.CheckDBFileExists(dbResultPath); errCheck != nil {
		logerr.FatalEvent("cannot check DB on exists", errCheck)
	}

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir(config.DefaultWebDir))
	logerr.InfoMsg("Running handlers.")
	r.Handle("/*", fileServer)

	r.Get("/api/nextdate", tp.NextDateHand)

	r.Get("/api/tasks", mwar.CheckAuth(tp.UpcomingTasksWithSearch))

	r.Post("/api/task/done", mwar.CheckAuth(tp.GeneratedNextDate))

	r.Post("/api/signin", tp.LoginAuth)

	r.Delete("/api/task", mwar.CheckAuth(tp.DeleteTaskScheduler))
	r.Get("/api/task", mwar.CheckAuth(tp.GetTaskByID))
	r.Post("/api/task", mwar.CheckAuth(tp.AddSchedulerTask))
	r.Put("/api/task", mwar.CheckAuth(tp.UpdateTask))

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
