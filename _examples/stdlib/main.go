package main

import (
	"fmt"
	"net/http"

	"github.com/lrstanley/mdhttp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Gopher!")
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/docs/", mdhttp.New("/docs/", http.Dir("../docs/"), nil))
	http.ListenAndServe(":8080", nil)
}
