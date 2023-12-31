package migrations

import (
	"gorm.io/gorm"
)

type Migration struct {
	DB *gorm.DB
}

func Migrate(db *gorm.DB) error {
	migrations := []func() error{}

	migration := Migration{
		DB: db,
	}

	migrations = append(
		migrations,
		migration.MigrateUsers,
		migration.MigrateTokens,
		migration.MigratePosts,
		migration.MigrateComments,
		migration.MigrateCommentReactions,
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
