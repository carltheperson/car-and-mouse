package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := ":8080"
	fmt.Println("Starting server on", port)
	http.ListenAndServe(port, http.FileServer(http.Dir(`./website`)))
}
