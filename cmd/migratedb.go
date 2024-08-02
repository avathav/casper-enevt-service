package cmd

import (
	"event-service/internal/database/gorm"
	"event-service/internal/di"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migratedb",
	Run: migrationHandler,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrationHandler(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()

	db, connectionErr := gorm.Connection(di.DatabaseParameters())
	if connectionErr != nil {
		log.WithContext(ctx).WithError(connectionErr).Panic("could not get db connection")
	}

	sdb, dbErr := db.DB()
	if dbErr != nil {
		log.WithContext(ctx).WithError(dbErr).Panic("could not get db connection")
	}

	driver, driverErr := mysql.WithInstance(sdb, &mysql.Config{})
	if driverErr != nil {
		log.WithContext(ctx).WithError(driverErr).Panic("could not get driver")
	}

	migration, migrationErr := migrate.NewWithDatabaseInstance(
		"file://migrations",
		di.MysqlDriver(), driver)
	if migrationErr != nil {
		log.WithContext(ctx).WithError(migrationErr).Panic("could not get driver")
	}

	if err := migration.Up(); err != nil {
		log.WithContext(ctx).WithError(err).Panic("could not process migrations")
	}
}
