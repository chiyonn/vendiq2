package main

import (
	"log"

	"github.com/chiyonn/vendiq2/pricer/internal/bot"
	"github.com/chiyonn/vendiq2/pricer/internal/client"
)

func main() {
	logger := log.Default()

	spapi, err := client.NewSPAPIClient(logger)
	if err != nil {
		panic(err)
	}

	pricerBot := bot.NewPricerBot(spapi)
	pricerBot.Run()
}
