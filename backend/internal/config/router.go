package config

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *App) NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*",
			"http://*",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", app.GetHi)

	//user route group
	r.Route("/user", func(r chi.Router) {

		r.Post("/", app.CreateClient)
		r.Delete("/", app.DeleteClient)
		r.Put("/{id}", app.UpdateClient)
	})
	//messanger route group
	r.Route("/messanger", func(r chi.Router) {
		r.Post("/", app.CreateMessanger)
		r.Put("/", app.UpdateMessanger)
		r.Get("/messages", app.GetMessages)
		r.Get("/allmessages", app.GetAllMessages)
	})

	r.Get("/usersbytag", app.GetUsersByTag)

	return r
}
