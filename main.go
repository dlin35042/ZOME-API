package main

import (
	"fmt"
	"net/http"
)

// corsMiddleware adds CORS headers to responses to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers to allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all domains, adjust as needed
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Add your handlers
	mux.HandleFunc("/airdrop/start", airdropHandler)

	// Wrap the mux with the CORS middleware
	handler := corsMiddleware(mux)

	// Start the server with CORS-enabled handler
	// Start the server
	fmt.Println("Starting server at port 8000...")
	http.ListenAndServe(":8080", handler)
}

func airdropHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "Starting Airdrop...")

	// Execute the airdrop logic
	// err := startAirdrop()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	fmt.Fprintln(w, "Airdrop initiated successfully.")
}
