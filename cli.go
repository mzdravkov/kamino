package main

import (
	"github.com/codegangsta/cli"
)

var App *cli.App

func init() {
	App = cli.NewApp()
	App.Name = "kamino"
	App.Usage = "A platform for creating distributed Application as a Service systems"
	App.Version = version
	App.Author = "Mihail Zdravkov"
	App.Email = "mihail0zdravkov@gmail.com"
	App.Flags = []cli.Flag{}
	App.Commands = []cli.Command{
		{
			Name:  "daemon",
			Usage: "Starts the kamino daemon. If no role flag is provided, by default the daemon will be a worker. Use the --sage flag to make it a sage node (router node). Note that a node can be both worker and sage at the same time if both flags are provided.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "sage, s",
					Usage: "Tells the kamino daemon that it's a sage node (router node).",
				},
				cli.BoolFlag{
					Name:  "worker, w",
					Usage: "Tells the kamino daemon that it's a worker node.",
				},
				cli.StringFlag{
					Name:  "consul, c",
					Value: "127.0.0.1:8500",
					Usage: "The address, on which the consul agent is listening for client access.",
				},
				cli.IntFlag{
					Name:  "kamino-serve-port, sp",
					Value: 3457,
					Usage: "The port on which the node will take requests for serving a tenant application.",
				},
				cli.IntFlag{
					Name:  "kamino-internal-port, ip",
					Value: 3458,
					Usage: "The port on which the node will recieve commands from other nodes.",
				},
			},
			Action: daemonizeIfNeeded,
		},
		{
			Name:  "deploy",
			Usage: "Deploys a new tenant on the system.",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port, p",
					Value: 0,
					Usage: "Port from the host system to which, the container's exposed port will be mapped. If you don't use this flag, Kamino will find a random free port.",
				},
			},
			Action: deploy,
		},
	}

	// Default action if no commands are provided
	App.Action = func(c *cli.Context) {
		if !c.Bool("daemon") {
			cli.ShowAppHelp(c)
		}
	}
}
