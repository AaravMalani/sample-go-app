package routes

import (
	"encoding/json"
	"net/http"

	"github.com/AaravMalani/cvwo-forum-backend/internal/api"
	"github.com/AaravMalani/cvwo-forum-backend/internal/auth"
	"github.com/AaravMalani/cvwo-forum-backend/internal/handlers/comments"
	"github.com/AaravMalani/cvwo-forum-backend/internal/handlers/posts"
	"github.com/AaravMalani/cvwo-forum-backend/internal/handlers/topics"
	"github.com/AaravMalani/cvwo-forum-backend/internal/handlers/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

const (
	ErrorEncounteredMessage = "There was an error encountered while handling this request."
)

func GenericHandler(handler func(w http.ResponseWriter, req *http.Request) (*api.Response, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		response, err := handler(w, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else if response.ErrorCode != 0 {
			w.WriteHeader(response.ErrorCode)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
func GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "DELETE"},
			AllowCredentials: true,
		}))
		r.Use(jwtauth.Verifier(auth.JwtTokenAuth))
		r.Use(auth.AuthMiddleware)
		r.Post("/login", GenericHandler(users.HandleLogin))
		r.Post("/register", GenericHandler(users.HandleRegister))
		r.Get("/topics/list", GenericHandler(topics.HandleList))
		r.Get("/posts/list", GenericHandler(posts.HandleList))
		r.Get("/comments/list", GenericHandler(comments.HandleList))

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Authenticator)
			r.Get("/user/moderators", GenericHandler(users.HandleModerators))
			r.Get("/user/banned", GenericHandler(users.HandleBanned))

			r.Post("/posts/create", GenericHandler(posts.HandleCreate))
			r.Post("/posts/edit", GenericHandler(posts.HandleEdit))
			r.Post("/posts/vote", GenericHandler(posts.HandleVote))
			r.Post("/comments/create", GenericHandler(posts.HandleCreate))
			r.Post("/comments/edit", GenericHandler(posts.HandleEdit))
			r.Post("/comments/vote", GenericHandler(posts.HandleVote))

			r.Post("/topics/create", GenericHandler(topics.HandleCreate))
		})
	}
}
