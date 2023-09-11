package storage

import (
	"fmt"

	"github.com/MarselBissengaliyev/ggp-blog/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DB_Host,
		config.DB_User,
		config.DB_Password,
		config.DB_Name,
		config.DB_Port,
		config.DB_SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	return db, nil
}
