package main

import (
	"log"
	"net/http"
)

const (
	portservice = ":8080"
)

func GetCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Headers...")
	// Loop over header names
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			log.Println(name, value)
		}
	}
	// get requested metod
	metod := r.Method
	log.Printf("The metod used %s", metod)
	// get requested host
	host := r.Host
	log.Println(host)

}

func main() {
	log.Printf("Application is listening port %s", portservice)
	http.HandleFunc("/health", GetCheck) // func to check service
	log.Fatal(http.ListenAndServe(portservice, nil))
}
