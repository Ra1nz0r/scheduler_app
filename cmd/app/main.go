package main

import (
	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"

	"github.com/ra1nz0r/scheduler_app/internal/server"
)

func main() {
	server.Run()
}
