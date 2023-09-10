package migrations

import "gorm.io/gorm"

type Migration struct {
	db *gorm.DB
}

func Migrate(db *gorm.DB) error {
	migration := Migration{
		db: db,
	}

	if err := migration.MigrateUsers(); err != nil {
		return err
	}

	return nil
}
