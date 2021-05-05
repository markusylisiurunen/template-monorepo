package migrations

import (
	"github.com/golang-migrate/migrate/v4"

	// https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(folder string, postgresURL string) error {
	sourceURL := "file://" + folder

	m, err := migrate.New(sourceURL, postgresURL)
	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}

		// TODO: what if the database is locked by another migration?

		return err
	}

	return nil
}
