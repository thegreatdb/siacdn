package cli

import (
	"os"

	urfavecli "github.com/urfave/cli"
)

func RunSiaCDNCLIApp() {
	// Create and run the app
	app := urfavecli.NewApp()
	app.Name = "SiaCDN"
	app.Usage = "Management binary for the SiaCDN backend"
	app.Commands = []urfavecli.Command{
		ServeCommand(),
	}
	app.Run(os.Args)
}