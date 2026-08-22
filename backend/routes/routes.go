package routes

import (
	"backend/auth"
	"backend/controllers"
	"backend/websocket"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func InitializeRoutes(c controllers.Controller) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-type", "withCredentials", "Access-Control-Allow-Origin : *", "Authorization"},
		AllowCredentials: true,
	}))

	r.HandleFunc("/ws/{userID}", websocket.HandleConnection)

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", c.RegisterUser)
		r.Post("/login", c.LoginUser)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.TokenAuth))
			r.Use(jwtauth.Authenticator(auth.TokenAuth))
			r.Get("/profile", c.Profile)
			r.Get("/users", c.GetAllUsers)
			r.Get("/getPublicKey", c.GetPublicKey)
		})
	})

	r.Route("/message", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.TokenAuth))
			r.Use(jwtauth.Authenticator(auth.TokenAuth))
			r.Post("/create", c.CreateMessage)
			r.Get("/view", c.GetMessages)
		})
	})

	r.Route("/friends", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.TokenAuth))
			r.Use(jwtauth.Authenticator(auth.TokenAuth))
			r.Get("/isFriend", c.IsFriend)
			r.Post("/sendFriendRequest", c.SendFriendRequest)
			r.Get("/getFriendRequests", c.GetFriendRequests)
			r.Get("/acceptFriendRequest", c.AcceptFriendRequest)
			r.Get("/rejectFriendRequest", c.RejectFriendRequest)
			r.Get("/isRequestSent", c.IsFriendRequestSent)
			r.Get("/isRequestReceived", c.IsRequestReceived)

		})
	})

	return r
}
