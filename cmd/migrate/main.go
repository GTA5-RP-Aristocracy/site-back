package main

import (
	"database/sql"
	"flag"
	"time"

	"github.com/GTA5-RP-Aristocracy/site-back/db/migrate"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func main() {
	// Read flags.
	path := flag.String("path", "", "path to the migration files")
	version := flag.String("version", "", "version of the migration (time format )")
	dbConnection := flag.String("db", "", "database connection string")

	flag.Parse()

	log.Info().
		Str("path", *path).
		Str("version", *version).
		Strs("flags", flag.Args()).
		Msg("")

	direction := flag.Arg(0)
	if direction == "" {
		log.Fatal().Msg("missing migration direction")
	}
	log.Info().Str("direction", direction).Msg("")

	dbInstance, err := sql.Open("postgres", *dbConnection)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to the database")
	}

	versionTime := time.Now()
	if *version != "" {
		versionTime, err = time.Parse(migrate.FormatVersion, *version)
		if err != nil {
			log.Fatal().Err(err).Msg("parse version")
		}
	}

	switch migrate.MigrationDirection(direction) {
	case migrate.MigrationUp:
		err = migrate.RunMigrationsUp(dbInstance, *path, versionTime)
	case migrate.MigrationDown:
		err = migrate.RunMigrationsDown(dbInstance, *path, versionTime)
	default:
		log.Fatal().Msg("invalid migration direction")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("run migrations")
	}

	log.Info().Msg("migrations run successfully")
}
