package main

import (
	"log"
	"net/http"

	"github.com/chiyonn/vendiq2/pricer/internal/db"
	"github.com/chiyonn/vendiq2/pricer/internal/router"
)

func main() {
	db.Init()
	db.Migrate()

	r := router.NewRouter()

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
