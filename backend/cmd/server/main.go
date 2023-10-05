package main

import (
	"log"
	"messanger/internal/config"
)

func main() {
	app := config.NewApp()
	app.InitiateMessaging()
	if err := app.Serv.ListenAndServe(); err != nil {
		log.Fatal("server couldn't run")
	}
}
