package helpers

import (
	"be-kreditkita/src/config"
	"be-kreditkita/src/models"
)

func Migrate() {
	config.DB.AutoMigrate(
		&models.Consumer{},
		&models.Limit{},
		&models.Transaction{},
	)
}
