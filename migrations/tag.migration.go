package migrations

import "github.com/MarselBissengaliyev/ggp-blog/models"

func (m *Migration) MigrateTags() error {
	if err := m.db.AutoMigrate(&models.Tag{}); err != nil {
		return err
	}

	return nil
}
