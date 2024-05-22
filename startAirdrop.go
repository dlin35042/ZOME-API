package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

func startAirdrop() error {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 1. Fetch the total supply of the token
	fmt.Println("------------- Airdrop Distribution --------------")
	supply := bscTokenSupply() // big.Int
	fmt.Printf("Total Supply: 			%s\n", supply.String())

	// 5% of supply will be used for airdrop
	tokensForAirdrop := new(big.Int).Div(new(big.Int).Mul(supply, big.NewInt(5)), big.NewInt(100))
	fmt.Printf("Tokens for Airdrop: 		%s\n", tokensForAirdrop.String())

	// 95% for Traders
	tokensForTraders := new(big.Int).Div(new(big.Int).Mul(tokensForAirdrop, big.NewInt(95)), big.NewInt(100))
	fmt.Printf("Tokens for Traders: 		%s\n", tokensForTraders.String())

	// 5% for pending order owners
	tokensForPendingOrders := new(big.Int).Div(new(big.Int).Mul(tokensForAirdrop, big.NewInt(5)), big.NewInt(100))
	fmt.Printf("Tokens for Pending Orders: 	%s\n", tokensForPendingOrders.String())

	// 2. Fetch the list of traders who have traded the token
	fmt.Println("------------- Fetch Pending Order Traders --------------")
	listingUsers := getPendingData()
	fmt.Println("Listings:", listingUsers)

	// 3. Distribute the tokens to the traders
	fmt.Println("------------- Distribute Tokens --------------")
	// Make recipient and amounts array from the listingUsers
	var recipients []common.Address
	var amounts []*big.Rat

	// Get sum of values of listingusers
	sum := new(big.Rat)

	for _, user := range listingUsers {
		value := new(big.Rat).SetFloat64(user.Value)
		sum.Add(sum, value)
	}

	fmt.Printf("Sum: %s\n", sum.FloatString(2))

	for _, user := range listingUsers {
		address := common.HexToAddress(user.Address) // Convert string to common.Address
		recipients = append(recipients, address)
		value := new(big.Rat).SetFloat64(user.Value)

		// Get percentage which value is of sum of  total users value
		percentage := new(big.Rat).Quo(value, sum)
		tokens := new(big.Rat).SetInt(tokensForPendingOrders) // Convert tokensForPendingOrders to *big.Rat
		amount := new(big.Rat).Mul(percentage, tokens)
		amounts = append(amounts, amount)
	}

	fmt.Println("Recipients:", recipients)
	fmt.Println("Amounts:", amounts)

	// distributeTokens(recipients, amounts)

	// ----------------- Fetch Trade Data  Distribution -----------------
	fmt.Println("------------- Fetch Trade Data --------------")
	traderUsers := getTradeData()
	fmt.Println("traderUsers:", traderUsers)

	// 3. Distribute the tokens to the traders
	fmt.Println("------------- Distribute Tokens --------------")
	// Make recipient and amounts array from the listingUsers
	var traderRecipients []common.Address
	var traderAmounts []*big.Rat

	// Get sum of values of listingusers
	traderSum := new(big.Rat)

	for _, user := range traderUsers {
		value := new(big.Rat).SetFloat64(user.Value)
		traderSum.Add(traderSum, value)
	}

	fmt.Printf("Sum: %s\n", traderSum.FloatString(2))

	for _, user := range traderUsers {
		address := common.HexToAddress(user.Address) // Convert string to common.Address
		traderRecipients = append(traderRecipients, address)
		value := new(big.Rat).SetFloat64(user.Value)

		// Get percentage which value is of sum of  total users value
		percentage := new(big.Rat).Quo(value, traderSum)
		tokens := new(big.Rat).SetInt(tokensForTraders) // Convert tokensForTraders to *big.Rat
		amount := new(big.Rat).Mul(percentage, tokens)
		traderAmounts = append(traderAmounts, amount)
	}

	fmt.Println("Recipients:", traderRecipients)
	fmt.Println("Amounts:", traderAmounts)

	// distributeTokens(traderRecipients, traderAmounts)

	return nil
}
