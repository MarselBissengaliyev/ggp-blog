package migrations

import "github.com/MarselBissengaliyev/ggp-blog/models"

func (m *Migration) MigratePosts() error {
	if err := m.db.AutoMigrate(&models.Post{}); err != nil {
		return err
	}

	return nil
}

func (m *Migration) MigratePostReactions() error {
	if err := m.db.AutoMigrate(&models.PostReaction{}); err != nil {
		return err
	}

	return nil
}
