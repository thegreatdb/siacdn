package cli

import (
	"log"
	"net/http"

	"github.com/thegreatdb/siacdn/siacdn-backend/server"
	urfavecli "github.com/urfave/cli"
)

func ServeCommand() urfavecli.Command {
	return urfavecli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Runs the SiaCDN server",
		Action:  serveCommand,
	}
}

func serveCommand(c *urfavecli.Context) error {
	srv, err := server.NewHTTPDServer()
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
