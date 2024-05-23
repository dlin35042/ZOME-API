package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// APIConfig holds the API credentials
type APIConfig struct {
	APIKey     string
	SecretKey  string
	Passphrase string
	Project    string // This applies only to WaaS APIs
}

var apiConfig = APIConfig{
	APIKey:     "",
	SecretKey:  "",
	Passphrase: "",
	Project:    "",
}

func preHash(timestamp, method, requestPath string, params url.Values) string {
	var queryString string
	if method == "GET" && params != nil {
		queryString = "?" + params.Encode()
	}
	return timestamp + method + requestPath + queryString
}

func sign(message, secretKey string) string {
	hmac := hmac.New(sha256.New, []byte(secretKey))
	hmac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

func createSignature(method, requestPath string, params url.Values) (string, string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	message := preHash(timestamp, method, requestPath, params)
	signature := sign(message, apiConfig.SecretKey)
	return signature, timestamp
}

func sendGetRequest(requestPath string, params url.Values) (map[string]interface{}, error) {
	signature, timestamp := createSignature("GET", requestPath, params)

	headers := map[string]string{
		"OK-ACCESS-KEY":        apiConfig.APIKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": apiConfig.Passphrase,
		"OK-ACCESS-PROJECT":    apiConfig.Project,
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.okx.com"+requestPath, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value[0])
	}
	req.URL.RawQuery = query.Encode()

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type kv struct {
	Address string
	Value   float64
}

func getPendingData() []kv {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ZOME_NFT_ADDRESS := os.Getenv("ZOME_NFT_ADDRESS")

	apiConfig = APIConfig{
		APIKey:     os.Getenv("OKX_API_KEY"),
		SecretKey:  os.Getenv("OKX_SECRET_KEY"),
		Passphrase: os.Getenv("OKX_PASSPHRASE"),
		Project:    "",
	}

	getRequestPath := "/api/v5/mktplace/nft/markets/listings"
	getParams := url.Values{
		"chain":             []string{"bsc"},
		"collectionAddress": []string{ZOME_NFT_ADDRESS},
	}

	data, err := sendGetRequest(getRequestPath, getParams)
	if err != nil {
		fmt.Printf("Error retrieving data: %v\n", err)
		return nil
	}

	dataMap := data["data"].(map[string]interface{})
	listings := dataMap["data"].([]interface{})
	// fmt.Println("Listings:", len(listings))

	makerPrice := make(map[string]float64)
	for _, listing := range listings {
		listingData := listing.(map[string]interface{})
		maker := listingData["maker"].(string)
		priceStr := listingData["price"].(string)
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			log.Fatalf("Failed to convert price to float: %v", err)
		}

		if _, ok := makerPrice[maker]; ok {
			makerPrice[maker] += price
		} else {
			makerPrice[maker] = price
		}

	}

	var ss []kv
	for k, v := range makerPrice {
		ss = append(ss, kv{k, v})
	}

	// Sort by value
	for i := 0; i < len(ss); i++ {
		for j := i + 1; j < len(ss); j++ {
			if ss[i].Value < ss[j].Value {
				ss[i], ss[j] = ss[j], ss[i]
			}
		}
	}

	// Return only 10% of highest makers
	return ss[:int(float64(len(ss))*0.1)]

}
