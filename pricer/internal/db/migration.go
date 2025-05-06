package db

import (
	"log"

	"github.com/chiyonn/vendiq2/pricer/internal/model"
)

func Migrate() {
	err := DB.AutoMigrate(
		&model.Pricing{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("DB Migration complete.")
}
