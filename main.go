package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// Set up the HTTP server routes

	http.HandleFunc("/api/airdrop/", airdropHandler)

	http.HandleFunc("/api/calculate", calculateAirdropHandler)

	// Start the server
	fmt.Println("Starting server at port 8000...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func airdropHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getAvailableAirdropAmountByAddress(w, r)
	} else if r.Method == "POST" {
		clainAirdropToken(w, r)
	} else {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
}

func calculateAirdropHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	err := startAirdrop()
	if err != nil {
		fmt.Println("Error:", err)

		return
	}
}

func getAvailableAirdropAmountByAddress(w http.ResponseWriter, r *http.Request) {

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

func clainAirdropToken(w http.ResponseWriter, r *http.Request) {

	// Get address value from the request
	address := r.URL.Path[len("/api/airdrop/"):]

	fmt.Println("Claim Airdrop...", address)

	// Execute the airdrop logic

	filename := "airdrop_data.json"

	amount, err := getValueByAddress(filename, address)
	if err != nil {
		fmt.Println("Error:", err)
		// return "Nothing to claim" message
		response := struct {
			Message string `json:"message"`
		}{
			Message: "We are sorry, there is nothing to claim for this address.",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Printf("Amount for address %s: %s\n", address, amount.FloatString(2))
	// return amount to the user in the response boy {amount: "123.45"}

	// Distribute Tokens to the address
	var recipients []common.Address
	var amounts []*big.Rat

	// Add the address and amount to the recipients and amounts slices
	recipients = append(recipients, common.HexToAddress(address))
	amounts = append(amounts, amount)

	fmt.Println("Distributing tokens...", recipients, amounts)

	// distributeTokens(recipients, amounts)

	// Create a response struct and encode it to JSON
	response := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("You have successfully claimed %s tokens. Balance will be available in a short time.", amount.FloatString(2)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
