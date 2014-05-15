package main

import (
	"fmt"
	"flag"
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"github.com/cactus/go-statsd-client/statsd"
	"time"
	"github.com/rakyll/globalconf"
)

func main() {
	log.Print("Starting sendgrid webhook endpoint...")
	// Config format

	statsdHost := flag.String("statsd_host", "", "")

	conf, err := globalconf.New("sendgridstatsd")
	conf.ParseAll()

	// first create a client
	client, err := statsd.New(*statsdHost, "test-client")
	// handle any errors
	if err != nil {
		log.Fatal(err)
	}
	// make sure to clean up
	defer client.Close()

	http.HandleFunc("/",  func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, client)
	})

	log.Print("done\n")
	http.ListenAndServe(":9090", nil)
}

func handler(w http.ResponseWriter, r *http.Request, client *statsd.Client) {
	if (r.Method != "POST") {
		log.Fatal("Expected POST request")
	}

	fmt.Fprintf(w, "Thanks!")

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var ec []map[string]interface{}

	if err := json.Unmarshal(buf.Bytes(), &ec); err != nil {
		log.Fatal(err)
	}
	for _, item := range ec {
		err := client.Inc(item["event"].(string), 1, 1.0)
		// handle any errors
		if err != nil {
			log.Fatalf("Error sending metric: %+v", err)
		}

		log.Printf("%s(%s): %s\n", item["email"], time.Unix(int64(item["timestamp"].(float64)), 0), item["event"])
	}
}
