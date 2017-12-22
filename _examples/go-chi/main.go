package main

import (
	"fmt"
	"net/http"

	"github.com/lrstanley/mdhttp"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Gopher!")
	})

	r.Mount("/docs/", mdhttp.New("/docs/", http.Dir("../docs/"), nil))
	http.ListenAndServe(":8080", r)
}
