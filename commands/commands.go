package commands

import (
	"github.com/aravinth2094/GoProxy/actions"
	"github.com/urfave/cli/v2"
)

var (
	Tunnel = &cli.Command{
		Name:  "tunnel",
		Usage: "Tunnel",
		Subcommands: []*cli.Command{
			{
				Name:   "server",
				Usage:  "Tunnel Server",
				Action: actions.TunnelServer,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "listen",
						Usage:    "Listen address",
						Required: true,
						EnvVars: []string{
							"TUNNEL_SERVER_LISTEN_ADDRESS",
						},
					},
					&cli.StringFlag{
						Name:     "target",
						Usage:    "Proxy server address",
						Required: true,
						EnvVars: []string{
							"PROXY_SERVER_TARGET_ADDRESS",
						},
					},
				},
			},
			{
				Name:   "client",
				Usage:  "Tunnel Client",
				Action: actions.TunnelClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "listen",
						Usage:    "Listen address",
						Required: true,
						EnvVars: []string{
							"TUNNEL_SERVER_LISTEN_ADDRESS",
						},
					},
					&cli.StringFlag{
						Name:     "target",
						Usage:    "Proxy server address",
						Required: true,
						EnvVars: []string{
							"PROXY_SERVER_TARGET_ADDRESS",
						},
					},
				},
			},
		},
	}

	Proxy = &cli.Command{
		Name:   "proxy",
		Usage:  "Proxy",
		Action: actions.ProxyServer,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "listen",
				Usage:    "Listen address",
				Required: true,
				EnvVars: []string{
					"PROXY_SERVER_LISTEN_ADDRESS",
				},
			},
		},
	}
)
