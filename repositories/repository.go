package repositories

import (
	"github.com/MarselBissengaliyev/ggp-blog/config"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
	Config *config.Config
}
