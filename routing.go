package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/codegangsta/cli"

	"golang.org/x/net/publicsuffix"
)

var tenants = map[string]string{"llama": "lvh.me:3000", "baba": "fb.com", "ah": "youtube.com"}

func reverseProxyToKaminoWorker(w http.ResponseWriter, req *http.Request) {
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

func startSageServer() {
	// the "/" pattern will match all requests
	http.HandleFunc("/", reverseProxyToKaminoWorker)

	err := http.ListenAndServe(":3456", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func reverseProxyToTenantApplication(w http.ResponseWriter, req *http.Request) {
}

// This is used for reverse proxying requests from a sage node to corresponding tenant
// and returning the response back to the sage node.
// There is another server that will handle commands from the sage nodes.
func startWorkerReverseProxy() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", reverseProxyToTenantApplication)

	server := http.Server{Addr: ":3457", Handler: serveMux}

	server.ListenAndServe()
}

func handleSageCommands(w http.ResponseWriter, req *http.Request) {
}

// This is used for recieving commands from the sage nodes.
// This is NOT for reverse proxying the request between a sage node and a tenant.
func startWorkerSageListener() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", handleSageCommands)

	server := http.Server{Addr: ":3458", Handler: serveMux}

	server.ListenAndServe()
}

func startWorkerServers() {
	go startWorkerReverseProxy()
	go startWorkerSageListener()
}

func startServer(c *cli.Context) {
	if c.Bool("sage") {
		startSageServer()
	} else {
		startWorkerServers()
	}
}
