package main

import (
	"log"
	"net/http"

	"github.com/chiyonn/vendiq2/pricer/internal/db"
	"github.com/chiyonn/vendiq2/pricer/internal/di"
	"github.com/chiyonn/vendiq2/pricer/internal/router"
)

func main() {

	db.Init()
	db.Migrate()
	database := db.GetDB()
    container := di.NewContainer(database)

	r := router.NewRouter(container)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
