package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Set up the HTTP server routes
	http.HandleFunc("/airdrop/start", airdropHandler)

	// Start the server
	fmt.Println("Starting server at port 8000...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func airdropHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "Starting Airdrop...")

	// Execute the airdrop logic
	err := startAirdrop()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Airdrop initiated successfully.")
}
