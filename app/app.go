package app

import (
	"github.com/aravinth2094/GoProxy/commands"
	"github.com/urfave/cli/v2"
)

func CreateApp() *cli.App {
	return &cli.App{
		Name:  "goproxy",
		Usage: "Secure proxy firewall to access private networks",
		Commands: []*cli.Command{
			commands.Tunnel,
			commands.Proxy,
		},
	}
}
