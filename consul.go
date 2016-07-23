package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func FindWorkerByTenant(consulAddress, tenant string) (string, string, string, error) {
	response, err := http.Get(fmt.Sprintf("http://%s/v1/catalog/service/%s", consulAddress, tenant))

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	type Response struct {
		Host, Address, ServicePort string
	}

	var result Response
	if err := json.Unmarshal(body, &result); err == nil {
		log.Fatal(err)
	}

	if result.Host == "" || result.Address == "" || result.ServicePort == "" {
		return "", "", "", errors.New("Can't find tenant")
	}

	return result.Host, result.Address, result.ServicePort, err
}
