package main

import (
	"backend/pkg/psql"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newApp() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	client := psql.PostgresClient{
		Host:     "localhost",
		User:     "postgres",
		Password: "mysecretpassword",
		Dbname:   "orders",
		Port:     "5432",
		Sslmode:  "disable",
	}
	db, err := psql.NewPostgres(client)
	if err != nil {
		log.Fatal(err)
	}

	_ = db

	log.Println("database connected")

	return r
}
