package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/abdulkarimogaji/invoGenius/db"
	"github.com/abdulkarimogaji/invoGenius/middleware"
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
		// Handle the error (e.g., return an internal server error)
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func StartServer() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /", healthHandler)

	stack := middleware.CreateStack(middleware.Logging)

	server := http.Server{
		Addr:    ":8000",
		Handler: stack(router),
	}
	log.Println("Starting server at port 8000")
	return server.ListenAndServe()
}
