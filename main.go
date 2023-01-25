package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	PARAM_URL      = "sgproxy-url"
	PARAM_USERNAME = "sgproxy-username"
	PARAM_PASSWORD = "sgproxy-password"
)

var (
	DEFAULT_URL      string
	DEFAULT_USERNAME string
	DEFAULT_PASSWORD string
)

func main() {
	bind := flag.String("bind", "127.0.0.1", "Address to listen on, default '127.0.0.1'")
	port := flag.Int("port", 8000, "Port to listen on, default 8443")
	tls_crt := flag.String("tls-cert", "", "Listen using HTTPS, path to TLS certificate to use")
	tls_key := flag.String("tls-key", "", "Listen using HTTPS, path to TLS key to use")
	default_url := flag.String("url", "https://www.google.com", fmt.Sprintf("Default URL to redirect to if '%s' query param not present", PARAM_URL))
	default_username := flag.String("username", "", fmt.Sprintf("Default HTTP Auth username to redirect to if '%s' query param not present", PARAM_USERNAME))
	default_password := flag.String("password", "", fmt.Sprintf("Default HTTP Auth password to redirect to if '%s' query param not present", PARAM_PASSWORD))
	flag.Parse()

	DEFAULT_URL = *default_url
	DEFAULT_USERNAME = *default_username
	DEFAULT_PASSWORD = *default_password

	// Start server
	address := fmt.Sprintf("%s:%d", *bind, *port)
	http.HandleFunc("/", proxyPass)
	if *tls_crt != "" && *tls_key != "" {
		fmt.Printf("Listenning on: https://%s\n", address)
		log.Fatal(http.ListenAndServeTLS(address, *tls_crt, *tls_key, nil))
	} else {
		fmt.Printf("Listenning on: http://%s\n", address)
		log.Fatal(http.ListenAndServe(address, nil))
	}
}

func proxyPass(res http.ResponseWriter, req *http.Request) {
	// Check if there are params to override defaults
	url_string := DEFAULT_URL
	if req.FormValue(PARAM_URL) != "" {
		url_string = req.FormValue(PARAM_URL)
	}
	username := DEFAULT_USERNAME
	if req.FormValue(PARAM_USERNAME) != "" {
		username = req.FormValue(PARAM_USERNAME)
	}
	password := DEFAULT_PASSWORD
	if req.FormValue(PARAM_PASSWORD) != "" {
		password = req.FormValue(PARAM_PASSWORD)
	}

	// Parse target URL
	url, _ := url.Parse(url_string)
	targetQuery := url.RawQuery

	// Use a custom Director
	director := func(req *http.Request) {
		// ----------------------------------------------
		// This part is identical to the default director
		// ----------------------------------------------
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(url, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		// --------------------
		// Custom director code
		// --------------------
		req.Host = url.Host

		// Add auth if needed
		if username != "" || password != "" {
			req.SetBasicAuth(username, password)
		}

		// Remove our query params
		q := url.Query()
		q.Del(PARAM_URL)
		q.Del(PARAM_USERNAME)
		q.Del(PARAM_PASSWORD)
		req.URL.RawQuery = q.Encode()
		fmt.Printf("Redirecting to: %s\n", url_string)
	}
	proxy := httputil.ReverseProxy{Director: director}

	proxy.ServeHTTP(res, req)
}
