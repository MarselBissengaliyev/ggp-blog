package migrations

import "github.com/MarselBissengaliyev/ggp-blog/models"

func (m *Migration) MigrateTokens() error {
	if err := m.db.AutoMigrate(&models.Token{}); err != nil {
		return err
	}

	return nil
}
