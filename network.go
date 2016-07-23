package main

import (
	"io"
	"log"
	"net/http"
)

func joinKaminoNetworkAsSage(joinAddress string) {
}

func joinKaminoNetworkAsWorker(joinAddress string) {
	response, err := http.Get(joinAddress)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		buffer := make([]byte, 0)
		if _, err := io.ReadFull(response.Body, buffer); err != nil {
			log.Fatal(err)
		}
		log.Fatalf("Tried to join the kamino network, but got response with status code %d. The body is: %s", response.StatusCode, string(buffer))
	}
}
