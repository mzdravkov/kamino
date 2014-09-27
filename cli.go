package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
)

var App *cli.App

func init() {
	App = cli.NewApp()
	App.Name = "kamino"
	App.Usage = "Platform for creating distributed Application as a Service systems"
	App.Version = Version
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
			Action: func(c *cli.Context) {
				if c.Args().First() != "" {
					port := uint16(c.Int("port"))
					if port == 0 {
						port = findFreePort()
					}
					if err := Deploy(c.Args().First(), port); err != nil {
						log.Fatal(err)
						return
					}
				} else {
					fmt.Println("You have to pass a name for the tenant as the first argument to deploy. Use 'kamino help deploy' for more info")
				}
			},
		},
	}

	App.Action = func(c *cli.Context) {
		if c.Bool("daemon") {
			println("Starting kamino as a daemon...")
		} else {
			cli.ShowAppHelp(c)
		}
	}
}
