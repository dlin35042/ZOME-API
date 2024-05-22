package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func distributeTokens(recipients []common.Address, amounts []*big.Rat) {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DISTRIBUTOR_PRIVATE_KEY := os.Getenv("DISTRIBUTOR_PRIVATE_KEY")
	INFURA_API_KEY := os.Getenv("INFURA_API_KEY")
	RPC_URL := os.Getenv("RPC_URL")
	AIRDROP_CONTRACT_ADDRESS := os.Getenv("AIRDROP_CONTRACT_ADDRESS")

	client, err := ethclient.Dial(RPC_URL + INFURA_API_KEY)

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fmt.Println("DISTRIBUTOR_PRIVATE_KEY", DISTRIBUTOR_PRIVATE_KEY)

	privateKey, err := crypto.HexToECDSA(DISTRIBUTOR_PRIVATE_KEY)
	if err != nil {
		log.Fatalf("Failed to create private key: %v", err)
	}

	// Read the ABI from a JSON file
	abiData, err := ioutil.ReadFile("./contracts/abis/airdrop.json")
	if err != nil {
		log.Fatalf("Failed to read ABI: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Create an instance of the contract
	contractAddressHex := common.HexToAddress(AIRDROP_CONTRACT_ADDRESS)
	contract := bind.NewBoundContract(contractAddressHex, parsedABI, client, client, client)

	// Prepare to send a transaction
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey)))
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	// auth.GasLimit = uint64(300000) // You might want to adjust this value based on the method being called

	// Convert *big.Rat to *big.Int
	var amountsInt []*big.Int
	for _, amount := range amounts {
		amountInt := new(big.Int)
		amountInt.SetString(amount.FloatString(0), 10)
		amountsInt = append(amountsInt, amountInt)
	}

	tx, err := contract.Transact(auth, "distributeTokens", recipients, amountsInt)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent! Tx Hash: %s\n", tx.Hash().Hex())
	// Wait for the transaction to be mined
	fmt.Println("Waiting for transaction to be mined...")
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for transaction: %v", err)
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Println("Transaction successfully mined!")
	} else {
		fmt.Println("Transaction failed to mine.")
	}

}
