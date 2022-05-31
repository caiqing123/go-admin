package cmd

import (
	"github.com/spf13/cobra"

	"api/database/seeders"
	"api/pkg/console"
	"api/pkg/database"
	"api/pkg/seed"
)

var CmdPlay = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database.",
	Run:   runPlay,
}

func runPlay(cmd *cobra.Command, args []string) {
	if !database.DB.Migrator().HasTable("jobs") {
		migrator().Up()
		seeders.Initialize()
		seed.RunAll()
		console.Success("Done seeding.")
	} else {
		console.Success("The database has already been initialized.")
	}
}
