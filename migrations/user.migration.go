package migrations

import (
	"github.com/MarselBissengaliyev/ggp-blog/models"
)

func (m *Migration) MigrateUsers() error {
	err := m.db.AutoMigrate(&models.User{})
	return err
}
