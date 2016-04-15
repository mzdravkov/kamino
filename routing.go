package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var tenants = map[string]string{"llama": "lvh.me:3000", "baba": "fb.com", "ah": "youtube.com"}

func handler(w http.ResponseWriter, req *http.Request) {
	tldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(req.Host)
	if err != nil {
		return
	}

	// take everything before the top level domain and the site name
	// i.e. the chain of the subdomains
	subdomainChain := strings.Split(req.Host, tldPlusOne)[0]

	subdomainChainSlice := strings.Split(subdomainChain, ".")
	// remove the last element, which is the empty string
	subdomainChainSlice = subdomainChainSlice[:len(subdomainChainSlice)-1]

	// if there is no subdomain
	if len(subdomainChainSlice) == 0 {
		return
	}

	// get the top level subdomain
	subdomain := subdomainChainSlice[len(subdomainChainSlice)-1]

	// probably we should ignore www if it's the only subdomain?
	// if len(subdomainChainSlice) == 1 && subdomain == "www" {
	// 	return
	// }

	newHost := tenants[subdomain]
	if (len(subdomainChainSlice)) > 1 {
		innerSubdomains := subdomainChainSlice[:len(subdomainChainSlice)-1]
		newHost = strings.Join(append(innerSubdomains, tenants[subdomain]), ".")
	}

	log.Println(fmt.Sprintf("Recieved a request to %s. Redirecting it to %s, based on a rule for the subdomain %s", req.Host, newHost, subdomain))

	newHostUrl, err := url.Parse(newHost)
	if err != nil {
		panic(err)
	}
	newHostUrl.Scheme = "http"

	proxy := httputil.NewSingleHostReverseProxy(newHostUrl)

	newHostUrl.Host = newHost
	req.URL = newHostUrl

	proxy.ServeHTTP(w, req)
}

func server() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":3456", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
