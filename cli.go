package main

import (
	"github.com/codegangsta/cli"
)

var App *cli.App

func init() {
	App = cli.NewApp()
	App.Name = "kamino"
	App.Usage = "Platform for creating distributed Application as a Service systems"
	App.Version = version
	App.Author = "Mihail Zdravkov"
	App.Email = "mihail0zdravkov@gmail.com"
	App.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "daemon, d",
			Usage: "starts the kamino daemon",
		},
	}
	App.Commands = []cli.Command{
		{
			Name:      "deploy",
			ShortName: "d",
			Usage:     "deploys a new tenant on the system",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port, p",
					Value: 0,
					Usage: "Port from the host system to which, the container's exposed port will be mapped. If you don't use this flag, Kamino will find a random free port",
				},
			},
			Action: deploy,
		},
	}

	// Default action if no commands are provided
	App.Action = func(c *cli.Context) {
		if c.Bool("daemon") {
			daemonize()
		} else {
			cli.ShowAppHelp(c)
		}
	}
}
