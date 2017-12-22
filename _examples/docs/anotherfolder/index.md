### Index.md

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/lrstanley/mdhttp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	md := mdhttp.New("/", http.Dir("../docs/"), mdhttp.DefaultRenderer)
	http.ListenAndServe(":8080", md)
}
```
