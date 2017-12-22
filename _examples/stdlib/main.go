package main

import (
	"fmt"
	"net/http"

	"github.com/lrstanley/mdhttp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Gopher!\n")
}

func main() {
	http.Handle("/docs/", mdhttp.New("/docs/", http.Dir("../docs/"), nil))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
