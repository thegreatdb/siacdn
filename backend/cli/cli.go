package cli

import (
	"net/http"
	"os"
	"time"

	urfavecli "github.com/urfave/cli"
)

const URLRoot = "http://localhost:9095"

var SiaCDNSecretKey string
var cliClient http.Client

func init() {
	SiaCDNSecretKey = os.Getenv("SIACDN_SECRET_KEY")
	cliClient = http.Client{Timeout: time.Second * 6}
}

func RunSiaCDNCLIApp() {
	// Create and run the app
	app := urfavecli.NewApp()
	app.Name = "SiaCDN"
	app.Usage = "Management binary for the SiaCDN backend"
	app.Commands = []urfavecli.Command{
		ServeCommand(),
		KubeCommand(),
		ProxyCommand(),
	}
	app.Run(os.Args)
}
