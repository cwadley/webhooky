package main

import (
	"bytes"
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

func main() {
	port, endpoint, sslCert, sslKey, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing args: %s", err)
		fmt.Println(usage)
		os.Exit(1)
	}

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Method:", r.Method)
		fmt.Println("URL:", r.URL)
		fmt.Println("Headers:")
		fmt.Println(r.Header)
		fmt.Println("Body:")
		theBody := new(bytes.Buffer)
		theBody.ReadFrom(r.Body)
		fmt.Println(theBody.String())
		fmt.Fprintf(w, "OK")
	})

	if sslKey == "" {
		err = http.ListenAndServe(":"+port, nil)
	} else {
		err = http.ListenAndServeTLS(":"+port, sslCert, sslKey, nil)
	}

	if err != nil {
		fmt.Println(err)
	}

}

func parseArgs(rawArgs []string) (string, string, string, string, error) {
	var (
		prt  = ""
		endp = "/"
		cert = ""
		key  = ""
	)

	for i := 0; i < len(rawArgs)-1; i++ {
		switch arg := rawArgs[i]; arg {
		case "--port":
			prt = rawArgs[i+1]
			i++
		case "--endpoint":
			endp = rawArgs[i+1]
			i++
		case "--sslCert":
			cert = rawArgs[i+1]
			i++
		case "--sslKey":
			key = rawArgs[i+1]
			i++
		default:
			fmt.Printf("Unrecognized switch: %s", arg)
		}
	}

	if prt == "" {
		return "", "", "", "", fmt.Errorf("The --port switch is required. Supply a valid port.")
	}
	if key != "" && cert == "" {
		return "", "", "", "", fmt.Errorf("--sslKey specified, but --sslCert not specified. Supply a valid certificate path.")
	}
	if key == "" && cert != "" {
		return "", "", "", "", fmt.Errorf("--sslCert specified, but --sslKey not specified. Supply a valid key path.")
	}

	return prt, endp, cert, key, nil
}
