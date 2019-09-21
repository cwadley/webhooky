package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

const usage = `
Usage: go run webhooky.go
	--port 		<port to listen on> (required)
	--endpoint 	<endpoint to listen to> (optional, / is default)
	--sslCert	<path to SSL cert> (optional)
	--sslKey 	<path to SSL key> (optional)
`

type WebhookyConfig struct {
	port     string
	endpoint string
	cert     string
	key      string
}

func main() {
	config, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing args: %s", err)
		fmt.Println(usage)
		os.Exit(1)
	}

	fmt.Printf("Webhooky starting on endpoint %s, port %s\n", config.endpoint, config.port)
	if config.key == "" {
		http.HandleFunc(config.endpoint, vomiter)
		err = http.ListenAndServe(":"+config.port, nil)
	} else {
		mux := http.NewServeMux()
		mux.HandleFunc(config.endpoint, vomiter)

		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		srv := &http.Server{
			Addr:         ":" + config.port,
			Handler:      mux,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		err = srv.ListenAndServeTLS(config.cert, config.key)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func vomiter(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	fmt.Println("Method:", r.Method)
	fmt.Println("URL:", r.URL)
	fmt.Println("Headers:")
	fmt.Println(r.Header)
	fmt.Println("Body:")
	theBody := new(bytes.Buffer)
	theBody.ReadFrom(r.Body)
	fmt.Println(theBody.String())
	fmt.Fprintf(w, "OK - Webhooky")
}

func parseArgs(rawArgs []string) (WebhookyConfig, error) {
	config := WebhookyConfig{endpoint: "/"}

	for i := 0; i < len(rawArgs)-1; i++ {
		switch arg := rawArgs[i]; arg {
		case "--port":
			config.port = rawArgs[i+1]
			i++
		case "--endpoint":
			config.endpoint = rawArgs[i+1]
			i++
		case "--sslCert":
			config.cert = rawArgs[i+1]
			i++
		case "--sslKey":
			config.key = rawArgs[i+1]
			i++
		default:
			fmt.Printf("Unrecognized switch: %s", arg)
		}
	}

	if config.port == "" {
		return WebhookyConfig{}, fmt.Errorf("The --port switch is required. Supply a valid port.")
	}
	if config.key != "" && config.cert == "" {
		return WebhookyConfig{}, fmt.Errorf("--sslKey specified, but --sslCert not specified. Supply a valid certificate path.")
	}
	if config.key == "" && config.cert != "" {
		return WebhookyConfig{}, fmt.Errorf("--sslCert specified, but --sslKey not specified. Supply a valid key path.")
	}

	return config, nil
}
