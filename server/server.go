package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/abdulkarimogaji/invoGenius/config"
	"github.com/abdulkarimogaji/invoGenius/db"
	"github.com/abdulkarimogaji/invoGenius/middleware"
	v1 "github.com/abdulkarimogaji/invoGenius/server/api/v1"
)

type HomePageData struct {
	Title   string
	Message string
}

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

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	data := HomePageData{
		Title:   "Invo Genius",
		Message: "Live",
	}

	// Render the template with data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func StartServer() error {
	router := http.NewServeMux()
	handler := v1.NewHandler()

	v1 := http.NewServeMux()
	router.HandleFunc("GET /health", healthHandler)
	router.HandleFunc("GET /home", homePageHandler)
	fs := http.FileServer(http.Dir("public"))
	router.Handle("/public/", http.StripPrefix("/public/", fs))
	router.Handle("/v1/api/", http.StripPrefix("/v1/api", v1))

	authRouter := http.NewServeMux()
	authRouter.HandleFunc("POST /users", handler.CreateUser)
	authRouter.HandleFunc("POST /invoices", handler.CreateInvoice)
	authRouter.HandleFunc("GET /invoices", handler.GetInvoices)
	authRouter.HandleFunc("GET /customers", handler.GetCustomers)
	authRouter.HandleFunc("GET /token", handler.CheckToken)
	v1.Handle("/", middleware.JwtAuthMiddleware(authRouter))

	v1.HandleFunc("POST /login", handler.Login)

	stack := middleware.CreateStack(middleware.Logging, middleware.AllowCORS)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", config.C.Port),
		Handler: stack(router),
	}

	log.Printf("Starting server at port %v", config.C.Port)

	return server.ListenAndServe()
}
