package migrations

import "github.com/MarselBissengaliyev/ggp-blog/models"

func (m *Migration) MigrateReports() error {
	if err := m.db.AutoMigrate(&models.Report{}); err != nil {
		return err
	}

	return nil
}
