package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/hello", helloHandle)
	log.Fatal(http.ListenAndServe("localhost:9999", nil))
}

func indexHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "WRL.Path = %q\n", req.URL.Path)
}

func helloHandle(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
