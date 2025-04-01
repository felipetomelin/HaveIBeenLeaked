package main

import (
	"SerasaLeaks/services/haveibeenleaked"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Importação do driver MySQL
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	address string
	db      *sql.DB
}

func NewApiServer(address string, db *sql.DB) *Server {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &Server{
		address: address,
		db:      db,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/haveibeenleaked").Subrouter()

	passwordStore := haveibeenleaked.NewStore(s.db)
	userHandler := haveibeenleaked.NewHandler(passwordStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.address)
	return http.ListenAndServe(s.address, router)
}
