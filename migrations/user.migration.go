package migrations

import (
	"github.com/MarselBissengaliyev/ggp-blog/models"
)

func (m *Migration) MigrateUsers() error {
	if err := m.DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}
