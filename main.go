package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("Registering routes....")
	router := httprouter.New()
	router.GET("/:place", resolve)

	fmt.Println("Starting the server...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func resolve(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	place := params.ByName("place")

	// Resolve from the list of places.
	config, err := ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	destination, ok := config.Places[place]
	if !ok {
		fmt.Fprintln(writer, "Could not find the destination for", place)
		return
	}

	http.Redirect(writer, request, destination, http.StatusTemporaryRedirect)
}

type Config struct {
	Places map[string]string
}

func ReadConfig() (*Config, error) {
	buffer, err := ioutil.ReadFile("places.yaml")
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(buffer, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
