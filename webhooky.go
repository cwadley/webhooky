package main

import (
	"fmt"
	"bytes"
	//"html"
	//"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	http.ListenAndServe(":8090", nil)
}
