package main

import (
	"log"
	"net/http"
)

func main() {
	/*
		// create file server handler
		fs := http.FileServer(http.Dir("/bundles"))

		// handle `/bundles` route
		http.Handle("/bundles/", http.StripPrefix("/bundles/", fs))
	*/

	// This works too, but "/static2/" fragment remains and need to be striped manually
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	// start HTTP server with `http.DefaultServeMux` handler
	log.Fatal(http.ListenAndServe(":9000", nil))

}
