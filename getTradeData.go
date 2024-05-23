package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type sv struct {
	Address string
	Value   float64
}

func getTradeData() []sv {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	OPENSEA_API_KEY := os.Getenv("OPENSEA_API_KEY")
	OPENSEA_NFT_COLLECTION_SLUG := os.Getenv("OPENSEA_NFT_COLLECTION_SLUG")
	AIRDROP_START := os.Getenv("AIRDROP_START")
	AIRDROP_END := os.Getenv("AIRDROP_END")
	event_type := "sale"

	url := fmt.Sprintf("https://api.opensea.io/api/v2/events/collection/%s?event_type=%s&after=%s&before=%s", OPENSEA_NFT_COLLECTION_SLUG, event_type, AIRDROP_START, AIRDROP_END)

	var allEvents []map[string]interface{}

	for {
		// fmt.Println("Fetching...", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return nil
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("x-api-key", OPENSEA_API_KEY)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return nil
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return nil
		}

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return nil
		}

		events, found := response["asset_events"].([]interface{})
		if !found {
			fmt.Println("No events found in response, ending pagination")
			break
		}

		// Print count of events fetched
		// fmt.Printf("Fetched %d events\n", len(events))

		for _, event := range events {
			allEvents = append(allEvents, event.(map[string]interface{}))
		}

		if len(events) >= 50 {
			// Check for the next cursor
			next, found := response["next"]
			if !found || next == nil {
				fmt.Println("No more pages.")
				break
			}

			url = fmt.Sprintf("https://api.opensea.io/api/v2/events/collection/%s?event_type=%s&after=%s&before=%s&next=%s", OPENSEA_NFT_COLLECTION_SLUG, event_type, AIRDROP_START, AIRDROP_END, next)

			time.Sleep(1 * time.Second) // Sleep to avoid hitting rate limit
		} else {
			// fmt.Println("No more events to fetch.")
			break
		}
	}

	// fmt.Printf("Total events fetched: %d\n", len(allEvents))

	// Get seller and sum of quantity of payment from all events
	sellerQuantity := make(map[string]float64)
	for _, event := range allEvents {
		seller := event["seller"].(string)
		quantity := event["payment"].(map[string]interface{})["quantity"].(string)
		quantityFloat, err := strconv.ParseFloat(quantity, 64)
		if err != nil {
			fmt.Println("Error converting quantity to float:", err)
			return nil
		}

		if _, ok := sellerQuantity[seller]; ok {
			sellerQuantity[seller] += quantityFloat
		} else {
			sellerQuantity[seller] = quantityFloat
		}
	}

	// Sort sellerQuantity by sum of quantity
	var listingUsers []sv
	for seller, quantity := range sellerQuantity {
		listingUsers = append(listingUsers, sv{Address: seller, Value: quantity})
	}

	return listingUsers

}
