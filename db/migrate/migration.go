package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	FormatVersion = "2006-01-02-15:04:05"

	MigrationUp   MigrationDirection = "up"
	MigrationDown MigrationDirection = "down"
)

type (
	MigrationDirection string
)

func CreateMigrationFiles(path string, name string) error {
	// Create the migration file.
	migrationFileUp, err := os.Create(
		filepath.Join(path, fmt.Sprintf("%s_%s.%s.sql", time.Now().Format(FormatVersion), name, MigrationUp)),
	)
	if err != nil {
		return fmt.Errorf("create migration file: %w", err)
	}
	defer migrationFileUp.Close()

	migrationFileDown, err := os.Create(
		filepath.Join(path, fmt.Sprintf("%s_%s.%s.sql", time.Now().Format(FormatVersion), name, MigrationDown)),
	)
	if err != nil {
		return fmt.Errorf("create migration file: %w", err)
	}
	defer migrationFileDown.Close()

	return nil
}

func RunMigrationsUp(db *sql.DB, path string, version time.Time) error {
	// Read the migration files.
	migrations, err := ReadMigrationFiles(path, version, MigrationUp)
	if err != nil {
		return fmt.Errorf("read migration files: %w", err)
	}

	// Run the migration files.
	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("run migration: %w", err)
		}
	}

	return nil
}

func RunMigrationsDown(db *sql.DB, path string, version time.Time) error {
	// Read the migration files.
	migrations, err := ReadMigrationFiles(path, version, MigrationDown)
	if err != nil {
		return fmt.Errorf("read migration files: %w", err)
	}

	// Run the migration files.
	for _, migration := range migrations {
		queries := strings.Split(migration, ";")
		for _, query := range queries {
			_, err := db.Exec(query)
			if err != nil {
				return fmt.Errorf("run migration: %w", err)
			}
		}
	}

	return nil
}

func ReadMigrationFiles(path string, version time.Time, direction MigrationDirection) ([]string, error) {
	var migrations []string

	if version == (time.Time{}) {
		version = time.Now()
	}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a migration file.
		if !info.IsDir() && strings.HasSuffix(info.Name(), string(direction)+".sql") {
			fileVersion, err := getVersionFromMigrationName(info.Name())
			if err != nil {
				return fmt.Errorf("get version from migration name: %w", err)
			}

			// Check if the migration file is older than the version.
			if direction == MigrationUp && fileVersion.After(version) ||
				direction == MigrationDown && fileVersion.Before(version) {
				return nil
			}

			// Add the migration file to the list.
			migrationFile, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
			if err != nil {
				return fmt.Errorf("open migration file: %w", err)
			}
			defer migrationFile.Close()

			// Read the migration file.
			migration := make([]byte, info.Size())
			_, err = migrationFile.Read(migration)
			if err != nil {
				return fmt.Errorf("read migration file: %w", err)
			}

			migrations = append(migrations, string(migration))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("traverse migration files: %w", err)
	}

	return migrations, nil
}

func getVersionFromMigrationName(migrationName string) (time.Time, error) {
	return time.Parse(FormatVersion, strings.Split(migrationName, "_")[0])
}
