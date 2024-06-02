package config

import "database/sql"

var DefaultWebDir = "./internal/web/"
var DefIPAddress = "0.0.0.0"
var DefaultPort = "7540"
var DbDefaultPath = "internal/storage_db/scheduler.db"

type DB struct {
	Db *sql.DB
}
