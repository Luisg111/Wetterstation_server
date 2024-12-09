package main

import (
	"luis/wetterserver/database"
	"luis/wetterserver/server"
)

func main() {
	var db = database.CreateNewSqliteDatabaseConnection()
	server.CreateNewHttpServer(&db)
}
