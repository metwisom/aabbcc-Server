package app

import (
	"aabbcc-Server/internal/handler/v1"
	"aabbcc-Server/internal/pkg/db"
	server2 "aabbcc-Server/internal/pkg/server"
	_ "github.com/mattn/go-sqlite3"
)

func Run() {

	db.Database.Connect()

	defer db.Database.Close()

	db.Database.CreateTable()

	var server server2.WebServer = server2.NewFiberWebServer()

	server.Post("/auth", v1.AuthPost)
	server.Post("/data", v1.DataPost)
	server.Get("/data", v1.DataGet)
	server.Listen(":45214")

}
