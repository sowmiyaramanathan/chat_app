package routes

import (
	"backend/auth"
	"backend/controllers"
	"backend/metrics"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

func InitializeRoutes(c controllers.Controller) *chi.Mux {
	r := chi.NewRouter()

	_ = godotenv.Load("../../.env")
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	r.Use(slogRequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Content-type", "Authorization"},
		AllowCredentials: true,
	}))

	r.HandleFunc("/ws/{userID}", c.HandleConnection)

	r.Get("/metrics", metrics.Handler)
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
			r.Use(metrics.HTTPMiddleware)
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
			r.Post("/sendRequest", c.SendFriendRequest)
			r.Get("/getFriendRequests", c.GetFriendRequests)
			r.Put("/acceptRequest", c.AcceptFriendRequest)
			r.Put("/rejectRequest", c.RejectFriendRequest)
			r.Get("/isRequestSent", c.IsFriendRequestSent)
			r.Get("/isRequestReceived", c.IsRequestReceived)

		})
	})

	return r
}

// slogRequestLogger records request metadata without exposing query parameters,
// which may contain credentials such as the WebSocket token.
func slogRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
