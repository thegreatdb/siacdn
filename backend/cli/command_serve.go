package cli

import (
	"log"
	"net/http"

	"github.com/thegreatdb/siacdn/backend/db"
	"github.com/thegreatdb/siacdn/backend/server"
	urfavecli "github.com/urfave/cli"
)

var databaseFilepath string

func ServeCommand() urfavecli.Command {
	return urfavecli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Runs the SiaCDN server",
		Action:  serveCommand,
		Flags: []urfavecli.Flag{
			urfavecli.StringFlag{
				Name:        "db",
				Value:       "database.json",
				Usage:       "Load database from `FILE` (default: 'database.json')",
				Destination: &databaseFilepath,
			},
		},
	}
}

func serveCommand(c *urfavecli.Context) error {
	database, err := db.OpenDatabase(databaseFilepath)
	if err != nil {
		log.Println("Could not open database: " + err.Error())
		return err
	}

	srv, err := server.NewHTTPDServer(database)
	if err != nil {
		log.Println("Could not create server: " + err.Error())
		return err
	}

	log.Println("Running server on port 9095")
	if err := http.ListenAndServe(":9095", srv); err != nil {
		log.Println("Error running server")
		return err
	}
	return nil
}
