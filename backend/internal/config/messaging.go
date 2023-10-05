package config

import (
	"context"
	"fmt"
	"log"
	"messanger/internal/entity"
	"time"
)

func (app *App) InitiateMessaging() {

	messangers, err := app.MRS.GetMessangers(context.Background())
	if err != nil {
		println(err)
	}

	for _, v := range messangers {
		err := app.InitiateFromMessangerOnReboot(v, context.Background())
		if err != nil {
			println(err)
		}
	}
}

func (app *App) InitiateFromMessanger(m entity.Messanger, ctx context.Context) error {
	clients, err := app.CS.GetClientsByTag(ctx, m.Filter)
	log.Println("InitiateFromMessanger", clients)
	if err != nil {
		return err
	}
	var messages []entity.Message
	for _, v := range clients {
		var message entity.Message
		message.Client = v
		message.Messanger = m
		message.TimeSent = time.Now()
		if err := app.MS.CreateMessage(ctx, &message); err != nil {
			return fmt.Errorf(fmt.Sprintf("from InitiateFromMessanger(): %s", err))
		}
		messages = append(messages, message)

	}
	err = app.MRS.MS.SendMessages(context.Background(), messages)
	if err != nil {
		println(err)
	}
	return nil
}

func (app *App) InitiateFromMessangerOnReboot(m entity.Messanger, ctx context.Context) error {
	t := time.Now()
	log.Println("InitiateFromMessangerOnReboot: ", m, "\nTime now: ", t)
	if m.TimeStart.Before(time.Now()) && m.TimeEnd.After(time.Now()) {
		go app.SendOnReboot(ctx, m)
	} else if m.TimeStart.After(time.Now()) {
		go func() {
			t := time.NewTimer(time.Until(m.TimeStart))
			<-t.C
			app.InitiateFromMessanger(m, ctx)
		}()

	}

	return nil
}

func (app *App) SendOnReboot(ctx context.Context, m entity.Messanger) error {
	clients, err := app.CS.GetClientsByTag(ctx, m.Filter)
	log.Println("InitiateFromMessanger", clients)
	if err != nil {
		return err
	}
	var messages []entity.Message
	messages, err = app.MS.GetUnsentMessages(ctx, m)
	if err != nil {
		return err
	}
	err = app.MRS.MS.SendMessages(context.Background(), messages)
	if err != nil {
		println(err)
	}
	return nil
}
