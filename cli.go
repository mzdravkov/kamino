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
			Usage: "Starts the kamino daemon. By default the daemon will be a worker. Use the --sage flag to make it a sage node (load balancer).",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "sage, s",
					Usage: "Tells the kamino daemon that it's a sage node (load balancer). If this flag is not present, the daemon will be a worker node.",
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
