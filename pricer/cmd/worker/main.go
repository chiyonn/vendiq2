package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/chiyonn/spapi/auth"
	"github.com/chiyonn/vendiq2/pricer/internal/bot"
	"github.com/chiyonn/vendiq2/pricer/internal/core"
)

func main() {
	logger := setupLogger()

	ctx := context.Background()
	if err := runBot(ctx ); err != nil {
		logger.Error("bot execution failed", slog.Any("err", err))
		os.Exit(1)
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func loadConfig() *auth.AuthConfig {
	return &auth.AuthConfig{
		RefreshToken: core.MustReadSecret("SPAPI_REFRESH_TOKEN"),
		ClientID:     core.MustReadSecret("LWA_CLIENT_ID"),
		ClientSecret: core.MustReadSecret("LWA_CLIENT_SECRET"),
	}
}

func runBot(ctx context.Context) error {
	cfg := loadConfig()
	b, err := bot.NewPricerBot(cfg)
	if err != nil {
		return err
	}
	
	if err := b.Run(ctx); err != nil {
		return err
	}
	return nil
}
