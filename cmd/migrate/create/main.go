package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/GTA5-RP-Aristocracy/site-back/db/migrate"
)

func main() {
	// Read flags.
	path := flag.String("path", "", "path to the migration files")
	// version := flag.String("version", "", "version of the migration")
	name := flag.String("name", "", "name of the migration")

	flag.Parse()

	// Create the migration files.
	if err := migrate.CreateMigrationFiles(*path, *name); err != nil {
		log.Fatal().Err(err).Msg("create migration files")
	}

	log.Info().Msg("migration files created")
}
