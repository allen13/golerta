package main

import (
	"fmt"
	"github.com/allen13/golerta/app"
	"github.com/allen13/golerta/app/auth/token"
	"github.com/allen13/golerta/app/config"
	"github.com/docopt/docopt-go"
	"log"
	"strings"
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
	configFile := args["--config"].(string)
	config := config.BuildConfig(configFile)

	if args["server"].(bool) {
		echo := app.BuildApp(config)
		log.Println("Starting golerta server...")

		var err error

		if config.Golerta.TLSEnabled {
			if config.Golerta.TLSAutoEnabled {
				tlsHosts := strings.Split(config.Golerta.TLSAutoHosts, ",")
				err = echo.StartAutoTLS(config.Golerta.BindAddr, tlsHosts, "cert-cache")
			} else {
				err = echo.StartTLS(config.Golerta.BindAddr, config.Golerta.TLSCert, config.Golerta.TLSKey)

			}
		} else {
			err = echo.Start(config.Golerta.BindAddr)
		}

		if err != nil {
			echo.Logger.Fatal(err)
		}

	}

	if args["createAgentToken"].(bool) {
		fmt.Println(token.CreateExpirationFreeAgentToken(args["<name>"].(string), config.Golerta.SigningKey))
	}
}
