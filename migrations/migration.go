package migrations

import (
	"gorm.io/gorm"
)

type Migration struct {
	db *gorm.DB
}

func Migrate(db *gorm.DB) error {
	migrations := []func() error{}

	migration := Migration{}

	migrations = append(
		migrations,
		migration.MigrateUsers,
		migration.MigrateTokens,
		migration.MigrateComments,
		migration.MigrateCommentReactions,
		migration.MigratePosts,
		migration.MigratePostReactions,
		migration.MigrateReports,
		migration.MigrateTags,
	)

	for _, m := range migrations {
		if err := m(); err != nil {
			return err
		}
	}

	return nil
}
