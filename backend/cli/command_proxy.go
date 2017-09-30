package cli

import (
	"log"
	"net/http"

	"github.com/thegreatdb/siacdn/backend/proxy"
	urfavecli "github.com/urfave/cli"
)

func ProxyCommand() urfavecli.Command {
	return urfavecli.Command{
		Name:    "proxy",
		Aliases: []string{"p"},
		Usage:   "Proxies http requests from a central handler to customer services",
		Action:  proxyCommand,
	}
}

func proxyCommand(c *urfavecli.Context) error {
	srv, err := proxy.New()
	if err != nil {
		log.Println("Could not create server: " + err.Error())
		return err
	}

	log.Println("Running server on port 9096")
	if err := http.ListenAndServe(":9096", srv); err != nil {
		log.Println("Error running server")
		return err
	}
	return nil
}
