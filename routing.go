package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"

	"golang.org/x/net/publicsuffix"
)

var tenants = map[string]string{"llama": "lvh.me:3000", "baba": "fb.com", "ah": "youtube.com"}

var consulAddress string

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

	_, address, port, err := FindWorkerByTenant(consulAddress, subdomain)
	if err != nil {
		log.Println(fmt.Sprintf("ERROR: The node recieved a request for tenant \"%s\", but cannot find a worker that hosts such tenant", subdomain))
		return
	}

	newHost := address + ":" + port
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

func startSageServer(servePort int) {
	log.Println("Sage node listening on port ", servePort)
	// the "/" pattern will match all requests
	http.HandleFunc("/", reverseProxyToKaminoWorker)

	err := http.ListenAndServe(":"+strconv.Itoa(servePort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func reverseProxyToTenantApplication(w http.ResponseWriter, req *http.Request) {
}

// This is used for reverse proxying requests from a sage node to a corresponding tenant
// and returning the response back to the sage node.
// There is another server that will handle commands from the sage nodes.
func startWorkerReverseProxy(servePort int) {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", reverseProxyToTenantApplication)

	server := http.Server{Addr: ":" + strconv.Itoa(servePort), Handler: serveMux}

	server.ListenAndServe()
}

func handlePingFromSage(w http.ResponseWriter, req *http.Request) {
	type PingResponse struct {
		Load float64
	}

	pingResponse := PingResponse{Load: getWorkload()}

	jsonResponse, err := json.Marshal(pingResponse)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := w.Write(jsonResponse); err != nil {
		log.Fatal(err)
	}
}

func handleDeployCommand(w http.ResponseWriter, req *http.Request) {
}

// This is used for recieving commands from the sage nodes.
// This is NOT for reverse proxying the request between a sage node and a tenant.
func startWorkerSageListener(internalPort int) {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", handlePingFromSage)
	serveMux.HandleFunc("/deploy", handleDeployCommand)

	server := http.Server{Addr: ":" + strconv.Itoa(internalPort), Handler: serveMux}

	server.ListenAndServe()
}

func startWorkerServers(servePort, internalPort int) {
	go startWorkerReverseProxy(servePort)
	go startWorkerSageListener(internalPort)
}

func startServer(context *cli.Context) {
	consulAddress = context.String("consul")
	if context.Bool("sage") {
		log.Println("ihaa ", 42)
		go startSageServer(context.Int("kamino-serve-port"))
	}
	if context.Bool("worker") {
		go startWorkerServers(context.Int("kamino-serve-port"), context.Int("kamino-internal-port"))
	}

	// make the main goroutine hang, so that it doesn't terminate
	select {}
}
