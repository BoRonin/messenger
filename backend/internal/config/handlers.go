package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"messanger/internal/entity"
	"messanger/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *App) GetHi(w http.ResponseWriter, r *http.Request) {
	m, _ := json.Marshal(map[string]string{
		"message": "hi",
	})
	w.Write([]byte(m))
}

// Client Handlers
func (app *App) CreateClient(w http.ResponseWriter, r *http.Request) {
	htc := entity.HttpToClient{}
	json.NewDecoder(r.Body).Decode(&htc)
	defer r.Body.Close()
	if !entity.IsTelValid(htc.Tel) {
		utils.ErrorJSON(w, errors.New("wrong number... must be of this pattern: 7xxxxxxxxxx"), http.StatusBadRequest)
		return
	}
	client := entity.NewClient(htc)
	if err := app.CS.CreateClient(context.Background(), &client); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, client)
}
func (app *App) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intid, err := strconv.Atoi(id)
	if err != nil {
		utils.ErrorJSON(w, errors.New("couldn't get id"), http.StatusBadRequest)
		return
	}
	if intid == 0 {
		utils.ErrorJSON(w, errors.New("bad request"), http.StatusBadRequest)
		return
	}
	if err := app.CS.DeleteClient(context.Background(), entity.Client{ID: uint(intid)}); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, "deleted")
}

func (app *App) GetUsersByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	users, err := app.CS.DB.GetClientsByTag(context.Background(), tag)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, users)

}
func (app *App) UpdateClient(w http.ResponseWriter, r *http.Request) {
	idstring := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idstring)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
	}
	htc := entity.HttpToClient{}
	if err := json.NewDecoder(r.Body).Decode(&htc); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()
	client := entity.NewClient(htc)
	client.ID = uint(id)
	if err := app.CS.UpdateClient(context.Background(), client); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, client)

}

// messanger handlers
func (app *App) CreateMessanger(w http.ResponseWriter, r *http.Request) {
	me := entity.Messanger{}
	if err := json.NewDecoder(r.Body).Decode(&me); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		r.Body.Close()
		return
	}
	fmt.Printf("me: %v\n", me)
	defer r.Body.Close()
	if err := app.MRS.CreateMessanger(context.Background(), &me); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	if me.TimeStart.Before(time.Now()) && me.TimeEnd.After(time.Now()) {
		go app.InitiateFromMessanger(me, context.Background())
	} else if me.TimeStart.After(time.Now()) {
		go func() {
			t := time.NewTimer(time.Until(me.TimeStart))
			<-t.C
			app.InitiateFromMessanger(me, context.Background())
		}()

	}
	utils.WriteJSON(w, me)
}

func (app *App) UpdateMessanger(w http.ResponseWriter, r *http.Request) {
	me := entity.Messanger{}
	if err := json.NewDecoder(r.Body).Decode(&me); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()
	if err := app.MRS.UpdateMessanger(context.Background(), me); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	app.CancellationChan <- me.ID

	utils.WriteJSON(w, "all good")
}

// Message handlers
func (app *App) GetMessages(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	if id == 0 {
		utils.ErrorJSON(w, errors.New("wrong id, bro"), http.StatusBadRequest)
	}
	messages, err := app.MS.GetMessagesById(context.Background(), uint(id))
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
	}
	utils.WriteJSON(w, messages)
}
func (app *App) GetAllMessages(w http.ResponseWriter, r *http.Request) {

	messages, err := app.MS.GetMessages(context.Background())
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
	}
	utils.WriteJSON(w, messages)
}
