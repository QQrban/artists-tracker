package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// Serve dynamic files
	http.HandleFunc("/artist/", ArtistHandler)
	http.HandleFunc("/autocomplete", autocompleteHandler)
	http.HandleFunc("/filter", filterHandler)
	http.HandleFunc("/", ArtistsHandler)
	// Start server
	port := 8080
	fmt.Println("Listening on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}
