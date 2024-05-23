package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

type AirdropList struct {
	Address string  `json:"address"`
	Value   big.Rat `json:"value"`
}

func startAirdrop() error {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	airdropList := []AirdropList{}

	// 1. Fetch the total supply of the token
	fmt.Println("********************** Airdrop Distribution ***********************")
	supply := bscTokenSupply()
	tokensForAirdrop := new(big.Int).Div(new(big.Int).Mul(supply, big.NewInt(5)), big.NewInt(100))
	tokensForPendingOrders := new(big.Int).Div(new(big.Int).Mul(tokensForAirdrop, big.NewInt(5)), big.NewInt(100))
	tokensForTraders := new(big.Int).Div(new(big.Int).Mul(tokensForAirdrop, big.NewInt(95)), big.NewInt(100))
	fmt.Printf("Total Supply: 				%s\n", supply.String())
	fmt.Printf("Tokens for Airdrop: 		%s\n", tokensForAirdrop.String())
	fmt.Printf("Tokens for Pending Orders: 	%s\n", tokensForPendingOrders.String())
	fmt.Printf("Tokens for Traders: 		%s\n", tokensForTraders.String())

	// 2. Fetch the list of traders who have traded the token
	fmt.Println("********************** Fetch Pending Order Traders **********************")
	listingUsers := getPendingData()

	sum := new(big.Rat)

	for _, user := range listingUsers {
		value := new(big.Rat).SetFloat64(user.Value)
		sum.Add(sum, value)
	}

	for _, user := range listingUsers {
		address := common.HexToAddress(user.Address) // Convert string to common.Address
		value := new(big.Rat).SetFloat64(user.Value)

		// Get percentage which value is of sum of  total users value
		percentage := new(big.Rat).Quo(value, sum)
		tokens := new(big.Rat).SetInt(tokensForPendingOrders) // Convert tokensForPendingOrders to *big.Rat
		amount := new(big.Rat).Mul(percentage, tokens)

		// Add the address and amount to the airdropList
		airdropList = append(airdropList, AirdropList{Address: address.String(), Value: *amount})
		fmt.Printf("++++++++++++ Address: %s, Amount: %s\n", address.String(), amount.FloatString(2))
	}

	// ********************** Fetch Trade Data  Distribution **********************
	fmt.Println("********************** Fetch Trade Data **********************")
	traderUsers := getTradeData()
	traderSum := new(big.Rat)

	for _, user := range traderUsers {
		value := new(big.Rat).SetFloat64(user.Value)
		traderSum.Add(traderSum, value)
	}

	for _, user := range traderUsers {
		address := common.HexToAddress(user.Address)
		value := new(big.Rat).SetFloat64(user.Value)

		// Get percentage which value is of sum of  total users value
		percentage := new(big.Rat).Quo(value, traderSum)
		tokens := new(big.Rat).SetInt(tokensForTraders) // Convert tokensForTraders to *big.Rat
		amount := new(big.Rat).Mul(percentage, tokens)

		// Add the address and amount to the airdropList, if the address is already in the list, then add the amount to the existing amount
		found := false
		for i, airdrop := range airdropList {
			if airdrop.Address == address.String() {
				airdropList[i].Value.Add(&airdropList[i].Value, amount)
				found = true
				break
			}
		}
		if !found {
			airdropList = append(airdropList, AirdropList{Address: address.String(), Value: *amount})
		}
		fmt.Printf("------------ Address: %s, Amount: %s\n", address.String(), amount.FloatString(2))
	}

	// Save airdropList to a file
	if err := saveAirdropListToFile(airdropList, "airdrop_data.json"); err != nil {
		log.Fatalf("Failed to save airdrop data: %v", err)
	}

	return nil
}

func saveAirdropListToFile(airdropList []AirdropList, filename string) error {
	// Marshal the data into JSON format
	jsonData, err := json.MarshalIndent(airdropList, "", "    ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file
	if err := ioutil.WriteFile(filename, jsonData, 0644); err != nil {
		return err
	}

	return nil
}
func getValueByAddress(filename, address string) (*big.Rat, error) {
	// Read the file
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Deserialize JSON into a slice of AirdropList
	var airdropList []AirdropList
	if err := json.Unmarshal(jsonData, &airdropList); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Search for the address
	for _, airdrop := range airdropList {
		if airdrop.Address == address {
			// Return the amount associated with the address
			return &airdrop.Value, nil
		}
	}

	// Return nil and an error if the address is not found
	return nil, fmt.Errorf("address not found")
}
