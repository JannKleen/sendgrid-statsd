package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"bytes"
)

func main() {
	fmt.Printf("Starting sendgrid webhook endpoint...\n")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9090", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("%s(%s): %s\n", item["email"], item["timestamp"], item["event"])
	}
}
