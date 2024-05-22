package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// ERC20 ABI JSON to interact with the contract
const tokenABI = `[{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

func bscTokenSupply() *big.Int {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	INFURA_API_KEY := os.Getenv("INFURA_API_KEY")
	RPC_URL := os.Getenv("RPC_URL")
	ZOME_TOKEN_ADDRESS := os.Getenv("ZOME_TOKEN_ADDRESS")

	client, err := ethclient.Dial(RPC_URL + INFURA_API_KEY)

	if err != nil {
		log.Fatal(err)
	}

	// Replace the below address with the token contract address
	tokenAddress := common.HexToAddress(ZOME_TOKEN_ADDRESS)
	contract, err := abi.JSON(strings.NewReader(tokenABI))
	if err != nil {
		log.Fatal(err)
	}

	// Creating a call message
	callMsg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: contract.Methods["totalSupply"].ID,
	}

	ctx := context.Background()
	result, err := client.CallContract(ctx, callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the result to a big.Int
	totalSupply := new(big.Int)
	totalSupply.SetBytes(result)

	return totalSupply
}
