package main

import (
	"golang.org/x/net/publicsuffix"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var tenants = map[string]string{"llama": "http://lvh.me:3000", "baba": "http://fb.com", "ah": "http://youtube.com"}

func handler(w http.ResponseWriter, req *http.Request) {
	tldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(req.Host)
	if err != nil {
		return
	}

	subdomainChain := strings.Split(req.Host, tldPlusOne)[0]

	subdomainChainSlice := strings.Split(subdomainChain, ".")

	// if there is no subdomain
	if len(subdomainChainSlice) < 2 {
		return
	}

	// get the top level subdomain
	subdomain := subdomainChainSlice[len(subdomainChainSlice)-2]

	if subdomain == "www" {
		return
	}

	host, err := url.Parse(tenants[subdomain])
	if err != nil {
		panic(err)
	}
	host.Scheme = "http"

	// host without the top level subdomain
	// (the one which specifies the remote host of the tenant)
	newHost := tldPlusOne
	if len(subdomainChainSlice) > 2 {
		subdomains := strings.Join(subdomainChainSlice, ".")
		newHost = strings.Join([]string{subdomains, newHost}, ".")
	}
	req.Host = newHost

	log.Println(strings.Join([]string{subdomain, " subdomain will proxy to ", tenants[subdomain]}, ""))

	proxy := httputil.NewSingleHostReverseProxy(host)

	proxy.ServeHTTP(w, req)
}

func server() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":3456", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
