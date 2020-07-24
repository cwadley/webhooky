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
	--port <port to listen on> (required)
	--endpoint <endpoint to listen to> (optional, / is default)
	--cert <path to SSL cert> (optional)
	--key <path to SSL key> (optional)
	--body <string containing a custom response body to be returned> (optional)
`

type webhooky struct {
	port     string
	endpoint string
	cert     string
	key      string
	body     string
}

func main() {
	webhooky := webhooky{}
	err := webhooky.parseArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing args: %s", err)
		fmt.Println(usage)
		os.Exit(1)
	}
	webhooky.serve()
}

func (w *webhooky) serve() {
	fmt.Printf("Webhooky starting on endpoint %s, port %s\n", w.endpoint, w.port)
	var err error
	if w.key == "" {
		http.HandleFunc(w.endpoint, w.vomiter)
		err = http.ListenAndServe(":"+w.port, nil)
	} else {
		mux := http.NewServeMux()
		mux.HandleFunc(w.endpoint, w.vomiter)

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
			Addr:         ":" + w.port,
			Handler:      mux,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		err = srv.ListenAndServeTLS(w.cert, w.key)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func (w *webhooky) vomiter(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	fmt.Println("RemoteAddr:", r.RemoteAddr)
	fmt.Println("Host:", r.Host)
	fmt.Println("Method:", r.Method)
	fmt.Println("URL:", r.URL)
	fmt.Println("RequestURI:", r.RequestURI)
	fmt.Println("Headers:")
	for _, header := range r.Header {
		fmt.Println("\t", header)
	}
	fmt.Println("Body:")
	theBody := new(bytes.Buffer)
	theBody.ReadFrom(r.Body)
	fmt.Println(theBody.String())
	fmt.Fprintf(writer, w.body)
}

func (w *webhooky) parseArgs(rawArgs []string) error {
	w.endpoint = "/"
	w.body = "OK - Webhooky"

	for i := 0; i < len(rawArgs)-1; i++ {
		switch arg := rawArgs[i]; arg {
		case "--port":
			w.port = rawArgs[i+1]
			i++
		case "--endpoint":
			w.endpoint = rawArgs[i+1]
			i++
		case "--cert":
			w.cert = rawArgs[i+1]
			i++
		case "--key":
			w.key = rawArgs[i+1]
			i++
		case "--body":
			w.body = rawArgs[i+1]
			i++
		default:
			fmt.Printf("Unrecognized switch: %s", arg)
		}
	}

	if w.port == "" {
		return fmt.Errorf("the --port switch is required. Supply a valid port")
	}
	if w.key != "" && w.cert == "" {
		return fmt.Errorf("--key specified, but --cert not specified. Supply a valid certificate path")
	}
	if w.key == "" && w.cert != "" {
		return fmt.Errorf("--cert specified, but --key not specified. Supply a valid key path")
	}

	return nil
}
