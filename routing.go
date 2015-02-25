package main

import (
	"log"
	"net/http"
	// "io"
	"strings"
	"golang.org/x/net/publicsuffix"
	"net/http/httputil"
	"net/url"
)

var tenants = map[string]string {"test": "https://google.com", "baba": "https://fb.com"}

func handler(w http.ResponseWriter, req *http.Request) {
	tldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(req.Host)
	if err != nil {
		return
	}

	subdomainChain := strings.Split(req.Host, tldPlusOne)[0]

	subdomainChainSlice := strings.Split(subdomainChain, ".")

	if len(subdomainChainSlice) < 2 {
		return
	}

	subdomain := subdomainChainSlice[len(subdomainChainSlice)-2]

	if subdomain == "www" {
		return
	}

	host, err := url.Parse(tenants[subdomain])
	if err != nil {
		panic(err)
	}

	log.Println(strings.Join([]string {subdomain, " subdomain will proxy to ", tenants[subdomain]}, ""))

	proxy := httputil.NewSingleHostReverseProxy(host)
	proxy.ServeHTTP(w, req)

	// io.WriteString(w, "llama")
}

func server() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// func startReverseProxy() {
	
// }

// func main() {
// 	remote, err := url.Parse("http://google.com")
// 	if err != nil {
// 		panic(err)
// 	}

// 	proxy := httputil.NewSingleHostReverseProxy(remote)
// 	http.HandleFunc("/", handler(proxy))
// 	err = http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.URL)
// 		p.ServeHTTP(w, r)
// 	}
// }
