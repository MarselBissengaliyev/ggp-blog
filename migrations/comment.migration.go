package migrations

import "github.com/MarselBissengaliyev/ggp-blog/models"

func (m *Migration) MigrateComments() error {
	if err := m.db.AutoMigrate(&models.Comment{}); err != nil {
		return err
	}

	return nil
}

func (m *Migration) MigrateCommentReactions() error {
	if err := m.db.AutoMigrate(&models.CommentReaction{}); err != nil {
		return err
	}

	return nil
}
