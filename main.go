package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/danoand/utils"
	"github.com/levigross/grequests"
)

var (
	cData    configData
	contJSON = map[string]string{"Content-Type": "application/json"}
	contXML  = map[string]string{"Content-Type": "application/xml"}
)

// Function roBody configures a RequestOptions object with the appropriate type of body (XML, JSON)
func roBody(inOp *grequests.RequestOptions, inRq *tRequests) (rErr error) {
	// Is this a POST?
	if inRq.Method != "POST" {
		rErr = errors.New("ERROR: Fn roBody - expecting a POST but received a GET request")
		return
	}

	// Configure a JSON or XML request body
	switch inRq.Typebody {
	case "application/json":
		inOp.Headers = contJSON
		inOp.JSON = inRq.Body
	case "application/xml":
		inOp.Headers = contXML
		inOp.XML = inRq.Body
	default:
		rErr = errors.New("ERROR: Fn roBody - did not receive an XML or JSON request object")
	}

	return
}

func main() {
	log.Println("Starting execution...")

	// Create a file to which the program will write
	oFile, oErr := os.Create("outfile.txt")
	if oErr != nil {
		log.Fatalf("ERROR: Error creating an output file see: %v\nEnding processing.", oErr)
	}
	defer oFile.Close()

	// Import the yaml file data
	lErr := loadRequests("data/calldata.toml")
	if lErr != nil {
		log.Fatalf("ERROR: Error loading TOML data see: %v\nEnding processing.", lErr)
	}

	// Iterate through the request data
	for i := 0; i < len(cData.Requests); i++ {
		log.Printf("INFO: Working on request #: %v\n", i+1)
		var madeReq bool

		// Construct the next request
		var myReq grequests.RequestOptions
		var myResp *grequests.Response

		// Assign any parameters
		myReq.Params = cData.Requests[i].Params
		// Determine the method
		myReq.Data = cData.Requests[i].Params
		switch cData.Requests[i].Method {
		case "GET":
			myResp, _ = grequests.Get(cData.Requests[i].URL, &myReq)
			madeReq = true
		case "POST":
			// Configure the ResponseOption for the applicable body and body type
			pErr := roBody(&myReq, &cData.Requests[i])
			if pErr != nil {
				log.Printf("Error: configuring the POST request body\n")
			} else {
				myResp, _ = grequests.Post(cData.Requests[i].URL, &myReq)
				madeReq = true
			}
		default:
			log.Printf("ERROR: Can't determine the request method")
		}

		// Print out the respsonse
		if madeReq == true {
			myRawRespStr, _, _ := utils.DumpResponse(myResp.RawResponse)

			// Format a string to be written to the output file
			oStr := fmt.Sprintf("\n\n******\n\nThe response to #: %v is:\n\n%v\n\n*****\n\n", i, myRawRespStr)

			// Write the output string
			_, wErr := oFile.WriteString(oStr)
			if wErr != nil {
				log.Printf("ERROR: Writing the response to the output file see: %v\n", wErr)
			}
		}
	}

	log.Println("Ending execution...")
}
