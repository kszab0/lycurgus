package main

import (
	"log"
	"net/http"
)

func main() {
	proxyAddress := ":8080"
	enableProxy := true
	blocker := NewBlocker(enableProxy)
	log.Fatal(http.ListenAndServe(proxyAddress, blocker))
}
