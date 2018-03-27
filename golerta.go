package main

import (
	"fmt"
	"log"

	"github.com/allen13/golerta/app"
	"github.com/allen13/golerta/app/auth/token"
	appconfig "github.com/allen13/golerta/app/config"
	"github.com/docopt/docopt-go"
)

const version = "Golerta 0.1.0"
const usage = `Golerta.

Usage:
	golerta server [--config=<config>]
	golerta createAgentToken <name> [--config=<config>]
	golerta --help
	golerta --version

Options:
  --config=<config>            The golerta config [default: ./golerta.toml].
  --help                       Show this screen.
  --version                    Show version.
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		log.Fatalln(err)
	}

	var config appconfig.GolertaConfig
	configFile := args["--config"].(string)
	if configFile != "./golerta.toml" {
		config = appconfig.BuildConfig(configFile)
	} else {
		config = appconfig.BuildConfig("golerta")
	}

	if args["server"].(bool) {
		echo := app.BuildApp(config)
		log.Println("Starting golerta server...")

		var err error

		if config.App.TLSEnabled {
			if config.App.TLSAutoEnabled {
				err = echo.StartAutoTLS(config.App.BindAddr)
			} else {
				err = echo.StartTLS(config.App.BindAddr, config.App.TLSCert, config.App.TLSKey)
			}
		} else {
			err = echo.Start(config.App.BindAddr)
		}

		if err != nil {
			echo.Logger.Fatal(err)
		}

	}

	if args["createAgentToken"].(bool) {
		fmt.Println(token.CreateExpirationFreeAgentToken(args["<name>"].(string), config.App.SigningKey))
	}
}
