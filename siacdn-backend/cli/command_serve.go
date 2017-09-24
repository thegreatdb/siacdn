package cli

import (
	"fmt"
	"log"
	"net/http"

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
	srv := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	/*
		srv, err := server.NewHTTPDServer(db)
		if err != nil {
			log.Println("Could not create server: " + err.Error())
			return err
		}
	*/
	log.Println("Running server on port 9095")
	if err := http.ListenAndServe(":9095", srv); err != nil {
		log.Println("Error running server")
		return err
	}
	return nil
}
