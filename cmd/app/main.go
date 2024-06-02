package main

import (
	_ "modernc.org/sqlite"

	"github.com/ra1nz0r/scheduler_app/internal/server"
)

func main() {
	server.Run()
}
