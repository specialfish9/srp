package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type UrlMapping map[string]*url.URL

type RProxy struct {
	mapping UrlMapping
	saver   Saver
}

func NewRProxy(mapping UrlMapping, saver Saver) *RProxy {
	return &RProxy{mapping, saver}
}

func (rp *RProxy) failInternal(rw http.ResponseWriter, cause error) {
	rp.fail(rw, cause, http.StatusInternalServerError)
}

func (rd *RProxy) fail(rw http.ResponseWriter, cause error, code int) {
	rw.WriteHeader(code)
	_, err := fmt.Fprint(rw, cause.Error())

	if err != nil {
		log.Println("error returning error: " + err.Error())
	}
}

func (rd *RProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Received request for: %s\n", req.Host)

	splitted := strings.Split(req.Host, ".")

	if len(splitted) == 0 {
		rd.fail(rw, fmt.Errorf("serpe: missing key"), http.StatusBadRequest)
		return
	}

	key := splitted[0]

	url, ok := rd.mapping[key]

	if !ok {
		rd.fail(rw, fmt.Errorf("url %s not found", key), http.StatusNotFound)
	}

	// set req Host, URL and Request URI to forward a request to the origin server
	req.Host = url.Host
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.RequestURI = ""

	// save the response from the origin server
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		rd.failInternal(rw, err)
		return
	}

	// copy headers
	for k, v := range response.Header {
		for _, x := range v {
			rw.Header().Add(k, x)
		}
	}

	rw.Header().Add("x-serpe", "true")

	rw.WriteHeader(response.StatusCode)
	io.Copy(rw, response.Body)

	// save
	if rd.saver != nil {
		rd.saver.Save(&RequestSaving{
			Url:    req.URL.String(),
			Sender: "sender url placeholder",
			Method: req.Method,
			Path:   req.URL.Path,
		})
	}
}
