package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/chiyonn/vendiq2/pricer/internal/bot"
	"github.com/chiyonn/vendiq2/pricer/internal/core"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/client"
)

func main() {
	logger := setupLogger()

	cfg := loadConfig()
	spapiClient, err := client.New(cfg, logger)
	if err != nil {
		logger.Error("client initialization failed", slog.Any("err", err))
		os.Exit(1)
	}

	ctx := context.Background()
	if err := runBot(ctx, spapiClient); err != nil {
		logger.Error("bot execution failed", slog.Any("err", err))
		os.Exit(1)
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func loadConfig() *client.Config {
	return &client.Config{
		RefreshToken: core.MustReadSecret("SPAPI_REFRESH_TOKEN"),
		ClientID:     core.MustReadSecret("LWA_CLIENT_ID"),
		ClientSecret: core.MustReadSecret("LWA_CLIENT_SECRET"),
	}
}

func runBot(ctx context.Context, c *client.Client) error {
	b := bot.NewPricerBot(c)
	if err := b.Run(ctx); err != nil {
		return err
	}
	return nil
}
