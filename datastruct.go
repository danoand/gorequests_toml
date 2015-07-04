package main

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type configData struct {
	Base     tbase
	Requests []tRequests
}

type tbase struct {
	URL string
}

type tRequests struct {
	URL      string
	Usebase  bool
	Method   string
	Typebody string
	Body     string
	Params   map[string]string
}

// Function loadRequests loads the requests data structures
func loadRequests(inFName string) (retErr error) {
	var fbytes []byte

	// Read the YAML file data into a byte array
	fbytes, retErr = ioutil.ReadFile(inFName)
	if retErr != nil {
		log.Fatalln("Fn loadRequests: Error reading data to initialize Customers see:", retErr)
	}

	// Unmarshall the YAML data into the data structure
	retErr = toml.Unmarshal(fbytes, &cData)
	if retErr != nil {
		log.Fatalln("Fn loadCustomers: Error parsing the Customer TOML data see:", retErr)
	}

	return
}
