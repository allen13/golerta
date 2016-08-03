package main

import (
  "github.com/docopt/docopt-go"
  "github.com/allen13/golerta/app"
)

const version = "Golerta 0.0.1"
const usage = `Golerta.

Usage:
  golerta --config=<config>
  golerta --help
  golerta --version

Options:
  --config=<config>            The golerta config [default: /etc/golerta/golerta.conf].
  --help                       Show this screen.
  --version                    Show version.
`

func main() {
  args, _ := docopt.Parse(usage, nil, true, version, false)
  configFile := args["--config"].(string)
  http := app.BuildApp(configFile)
  http.Listen(":5608")
}
