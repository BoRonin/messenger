package config

import (
	"fmt"
	"messanger/infrastructure/postgres"
	"messanger/internal/client"
	repository "messanger/internal/repository/postgres"
	"messanger/internal/service"
	"net/http"
	"os"
)

const (
	webPort = "3000"
	PostUrl = "https://probe.fbrq.cloud/v1/send/"
	Token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYzMTMxMjgsImlzcyI6ImZhYnJpcXVlIiwibmFtZSI6Imh0dHBzOi8vdC5tZS9ZQm9yb25pbiJ9.bI_-mlqLcH0C_wF_JCVJnfpLKapKdSF7ol6QXQUSYFE"
)

type App struct {
	Serv             http.Server
	CS               service.ClientService
	MRS              service.MessangerService
	MS               service.MessageService
	CancellationChan chan uint
}

func NewApp() *App {

	//setting DB connection
	dsn := os.Getenv("DSN")
	DB := postgres.NewDB(dsn)

	//setting repositories
	DBC := repository.NewDBClient(DB)
	DBMS := repository.NewDBMessanger(DB)
	DBM := repository.NewDBMessage(DB)

	//setting services
	CS := service.NewClientService(DBC)
	MS := service.NewMessageService(DBM)
	MC := client.NewMClient(PostUrl, Token, MS)
	MRS := service.NewMessangerService(DBMS, MC)

	//making an app instance
	app := &App{
		CS:  *CS,
		MRS: *MRS,
		MS:  *MS,
		Serv: http.Server{
			Addr: fmt.Sprintf(":%s", webPort),
		},
		CancellationChan: make(chan uint),
	}

	app.Serv.Handler = app.NewRouter()

	return app
}
