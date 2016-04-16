package main

import (
	"log"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

// A function to get the average workload of the host
// Note: it only works for linux
func getWorkload() float64 {
	loadAvg, err := linuxproc.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Fatal("Couldn't read /proc/loadavg")
	}

	return loadAvg.Last15Min
}
