package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ksdme/goplaces/pkg/config"
)

// The location to load the configuration from.
var location string = "places.yaml"

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

	log.Println("Starting the server")
	log.Fatal(http.ListenAndServe(":8080", router))
}
