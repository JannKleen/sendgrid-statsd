package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"github.com/cactus/go-statsd-client/statsd"
	"time"
)

func main() {
	log.Print("Starting sendgrid webhook endpoint...")
	// first create a client
	client, err := statsd.New("127.0.0.1:8125", "test-client")
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
	fmt.Fprintf(w, "Thanks!")

	if (r.Method != "POST") {
		log.Fatal("Expected POST request")
		// should i return a 200 or sth else (which would stop sendgrid sending any more events)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var ec []map[string]interface{}

	if err := json.Unmarshal(buf.Bytes(), &ec); err != nil {
		log.Fatal(err)
	}
	for _, item := range ec {
		// Send a stat
		err := client.Inc("stat1", 42, 1.0)
		// handle any errors
		if err != nil {
			log.Fatalf("Error sending metric: %+v", err)
		}

		log.Printf("%s(%s): %s\n", item["email"], time.Unix(int64(item["timestamp"].(float64)), 0), item["event"])
	}
}
