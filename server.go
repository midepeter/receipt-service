package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	PORT          string
	ADDRESS       string
	UsersDB       *Database
	TransactionDB *Database
	log           *slog.Logger
}

func NewServer(userDB *Database, transactionDB *Database, log *slog.Logger) *Server {
	return &Server{
		PORT:          "4002",
		ADDRESS:       "localhost",
		UsersDB:       userDB,
		TransactionDB: transactionDB,
		log:           log,
	}
}

func (s *Server) Start() {
	log.Printf("Starting server on %s:%s\n", s.ADDRESS, s.PORT)

	srv, err := setupHTTPServer(slog.Default(), s.PORT, s)
	if err != nil {
		log.Fatalf("Failed to set up server: %v", err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Printf("Server started on %s:%s\n", s.ADDRESS, s.PORT)

	// Block forever to keep the server running
	select {}
}

func (s *Server) Shutdown() {
	log.Println("Shutting down server...")
	// Implement graceful shutdown logic here
	// For example, close database connections, stop background jobs, etc.
	// This is a placeholder for actual shutdown logic
	// log.Println("Server shutdown complete.")
	// Implement any cleanup logic here
	// For example, close database connections, stop background jobs, etc.
	// This is a placeholder for actual shutdown logic
	// log.Println("Server shutdown complete.")
}

func setupHTTPServer(logger *slog.Logger, port string, srv *Server) (*http.Server, error) {
	handler := srv.router()
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}, nil
}

func (s Server) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "API is healthy")
	})

	mux.HandleFunc("/api/generate/receipt", s.getTransactionReceipt)
	mux.HandleFunc("/api/generate/summary", s.getTransactionSummary)
	return mux
}
