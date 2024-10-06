package app

import (
	"aabbcc-Server/internal/handler/v1"
	"aabbcc-Server/internal/pkg/db"
	server2 "aabbcc-Server/internal/pkg/server"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/signal"
	"syscall"
)

func Run() {

	db.Database.Connect()

	defer db.Database.Close()

	db.Database.CreateTable()

	var server server2.WebServer = server2.NewFiberWebServer()
	server.Post("/auth", v1.AuthPost)
	server.Post("/register", v1.RegisterPost)
	server.Post("/data", v1.DataPost)
	server.Get("/data", v1.DataGet)

	go func() {
		server.Listen(":45214")
		fmt.Println("Server is Run")
	}()
	defer server.Close()

	// Ожидание сигнала для остановки
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Блокировка до получения сигнала
	<-stopChan
	server.Close()

	// После получения сигнала можно выполнить необходимые действия по остановке
	fmt.Println("Shutting down server...")

}
