package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ksdme/goplaces/pkg/config"
)

// The location to load the configuration from.
var port int
var host string
var location string

// The configuration that will be used for resolving places.
var configuration *config.Config

// Reload the configuration when the user demands it.
func reload(writer http.ResponseWriter, reader *http.Request) {
	log.Println("Reloading configuration")

	if cfg, err := config.LoadConfig(location); err == nil {
		configuration = cfg
		fmt.Fprintln(writer, "Configuration reloaded.")
	} else {
		log.Println("Could not reload configuration", err)
		fmt.Fprintln(writer, "Could not reload configuration.")
	}
}

// Route to redirect the user to the place destination.
func resolve(writer http.ResponseWriter, request *http.Request) {
	place := chi.URLParam(request, "place")
	if len(place) == 0 {
		fmt.Println(writer, "Missing place")
		return
	}

	destination, ok := configuration.Places[place]
	if !ok {
		fmt.Fprintln(writer, "Could not find the destination for go/", place)
		return
	}

	http.Redirect(writer, request, destination, http.StatusTemporaryRedirect)
}

func main() {
	// The command line interface for the application.
	flag.StringVar(&location, "location", "places.yaml", "The location of the configuration file.")
	flag.StringVar(&host, "host", "127.0.0.1", "The host to bind the server on.")
	flag.IntVar(&port, "port", 8000, "The port to bind the server on.")
	flag.Parse()

	log.Println("Loading configuration")
	if cfg, err := config.LoadConfig(location); err == nil {
		configuration = cfg
	} else {
		log.Fatalln("Could not load configuration", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	log.Println("Registering routes")
	router.Get("/reload", reload)
	router.Get("/{place}", resolve)

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
