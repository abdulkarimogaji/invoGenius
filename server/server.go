package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abdulkarimogaji/invoGenius/config"
	"github.com/abdulkarimogaji/invoGenius/db"
	"github.com/abdulkarimogaji/invoGenius/middleware"
	v1 "github.com/abdulkarimogaji/invoGenius/server/api/v1"
)

type healthResponse struct {
	Error      bool       `json:"error"`
	Message    string     `json:"message"`
	ServerTime *time.Time `json:"server_time,omitempty"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	serverTime, err := db.DB.PingDB(ctx)

	var response healthResponse

	if err != nil {
		response.Error = true
		response.Message = err.Error()

	} else {
		response.Error = false
		response.Message = "Ping successful"
		response.ServerTime = &serverTime
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Error {
		w.WriteHeader(http.StatusInternalServerError)
	} else {

		w.WriteHeader(http.StatusOK)
	}
	w.Write(jsonResponse)
}

func StartServer() error {
	router := http.NewServeMux()
	handler := v1.NewHandler()

	v1 := http.NewServeMux()
	v1.HandleFunc("GET /health", healthHandler)
	v1.Handle("/v1/api/", http.StripPrefix("/v1/api", router))

	authRouter := http.NewServeMux()
	authRouter.HandleFunc("GET /require-auth", healthHandler)
	authRouter.HandleFunc("POST /users", handler.CreateUser)
	router.Handle("/", middleware.JwtAuthMiddleware(authRouter))

	router.HandleFunc("POST /login", handler.Login)

	stack := middleware.CreateStack(middleware.Logging, middleware.AllowCORS)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", config.C.Port),
		Handler: stack(v1),
	}
	log.Printf("Starting server at port %v", config.C.Port)
	return server.ListenAndServe()
}
