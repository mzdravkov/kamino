package main

import (
	"os"
)

const version = "0.0.1"

func main() {
	if err := isItDaemon(); err != nil {
		panic(err)
	}
	App.Run(os.Args)
}
