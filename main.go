package main

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	os.Args = os.Args[1:]
}

func main() {
	command := os.Args[0]
	switch command {
	case "deploy":
		if err := parseDeploy(); err != nil {
			parseHelp()
			panic("An error has occured while trying to deploy. Make sure you have rights to call docker and use nginx, and both are running")
		}
	case "help":
		parseHelp()
	}
}

func parseDeploy() error {
	nameF := flag.String("name", "", "Name of the tenant to deploy. It will be the same as the created docker container's name")
	portF := flag.Int("port", int(findFreePort()), "Port of the host system to which, the container's exposed port will be maped. If you don't pass port, Kamino will find random free port")

	flag.Parse()

	return Deploy(*nameF, uint16(*portF))
}

func parseHelp() {
	fmt.Println("Kamino is a tool for easy deployment of web applications")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("deploy :deploys an application with the name passed by -name=<name> ")
	fmt.Println("       :  -name=<name> is the name of the tenant. Will be the same as the container's name")
	fmt.Println("       :  -port=<port> gives the port of the host system to which,")
	fmt.Println("       :   the container's exposed port will be maped. If you don't pass port,")
	fmt.Println("       :   Kamino will find random free port")
	fmt.Println("help   :prints this message :)")
	fmt.Println("")
}
