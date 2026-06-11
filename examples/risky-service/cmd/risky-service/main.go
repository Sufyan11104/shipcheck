package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "risky-service")
	})

	_ = http.ListenAndServe(":8080", nil)
}
