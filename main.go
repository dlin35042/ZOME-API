package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Set up the HTTP server routes

	// api/airdrop/0x1234567890   0x1234567890 is the address
	http.HandleFunc("/api/airdrop/", airdropHandler)

	// Start the server
	fmt.Println("Starting server at port 8000...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func airdropHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	// Get address value from the request
	address := r.URL.Path[len("/api/airdrop/"):]

	fmt.Println("Starting Airdrop...", address)

	// Execute the airdrop logic

	filename := "airdrop_data.json"

	amount, err := getValueByAddress(filename, address)
	if err != nil {
		fmt.Println("Error:", err)
		// return amount 0
		response := struct {
			Amount string `json:"amount"`
		}{
			Amount: "0.00",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Printf("Amount for address %s: %s\n", address, amount.FloatString(2))
	// return amount to the user in the response boy {amount: "123.45"}

	// Create a response struct and encode it to JSON
	response := struct {
		Amount string `json:"amount"`
	}{
		Amount: amount.FloatString(2), // Convert big.Rat to string with 2 decimal places
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
