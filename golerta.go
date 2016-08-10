package main

import (
	"github.com/allen13/golerta/app"
	"github.com/docopt/docopt-go"
	"log"
	"fmt"
	"github.com/allen13/golerta/app/auth/token"
	"github.com/allen13/golerta/app/config"
)

const version = "Golerta 0.0.1"
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

	if args["server"].(bool){
		http := app.BuildApp(config)
		http.Listen(":5608")
	}

	if args["createAgentToken"].(bool){
		fmt.Println(token.CreateExpirationFreeAgentToken(args["<name>"].(string), config.Golerta.SigningKey))
	}
}
