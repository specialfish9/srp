package main

import (
	"log"
	"net/http"
	"net/url"
)

func mustMakeUrl(s string) *url.URL {
	url, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return url
}

func main() {

	cfg := LoadConfig()

	mapping := make(map[string]*url.URL)

	for _, m := range cfg.Mappings {
		mapping[m.Key] = mustMakeUrl(m.Destination)
	}

	proxy := NewRProxy(mapping)

	handler := http.NewServeMux()
	handler.Handle("/", proxy)

	log.Print("Listining on port 8080...")
	http.ListenAndServe(":8080", handler)
}
