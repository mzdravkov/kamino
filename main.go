package main

import (
	"os"
)

const version = "0.0.1"

func main() {
	isItDaemon, err := isItDaemon()
	if err != nil {
		panic(err)
	}

	if isItDaemon {
		server()
	} else {
		App.Run(os.Args)
	}
}
