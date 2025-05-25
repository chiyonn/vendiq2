package main

import (
	"log"
	"net/http"

	"github.com/chiyonn/vendiq2/pricer/internal/core"
	"github.com/chiyonn/vendiq2/pricer/internal/db"
	"github.com/chiyonn/vendiq2/pricer/internal/bot"
	"github.com/chiyonn/vendiq2/pricer/internal/router"
	"github.com/chiyonn/vendiq2/pricer/internal/consumer"

	"github.com/chiyonn/spapi/auth"
)

func main() {
	db.Init()
	db.Migrate()

	cfg := &auth.AuthConfig{
		ClientID:     core.MustReadSecret("LWA_CLIENT_ID"),
		ClientSecret: core.MustReadSecret("LWA_CLIENT_SECRET"),
		RefreshToken: core.MustReadSecret("SPAPI_REFRESH_TOKEN"),
	}
	httpClient := &http.Client{}

	b, err := bot.NewPricerBot(cfg, httpClient)
	if err != nil {
		log.Fatalf("failed to create PricerBot: %v", err)
	}

	go consumer.StartConsumer(b)

	r := router.NewRouter()
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
