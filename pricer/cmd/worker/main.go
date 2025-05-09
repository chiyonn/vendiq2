package main

import (
	"context"
	"log"

	"github.com/chiyonn/vendiq2/pricer/internal/bot"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/types"
)

func main() {
	logger := log.Default()
	cfg := &types.Config{ /* 認証・URLなど */ }

	client := spapi.New(cfg, logger)
	bot := bot.NewPricerBot(client)

	ctx := context.Background()
	if err := bot.Run(ctx); err != nil {
		log.Fatalf("Bot failed: %v", err)
	}
}
