package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	apiKey    = ""
	secretKey = ""
	baseURL   = "https://api.binance.th"
)

func main() {
	// Set the request parameters
	params := map[string]string{
		"symbol":      "BTCTHB",
		"side":        "BUY",
		"type":        "LIMIT",
		"timeInForce": "GTC",
		"quantity":    "0.0002",
		"price":       "2000000",
	}

	// Send the request
	resp, err := sendRequest("POST", "/api/v1/order", params)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println(string(body))
}

func sendRequest(method, path string, params map[string]string) (*http.Response, error) {
	// Create a new HTTP client
	client := http.Client{}
	// Add timestamp to request params
	params["timestamp"] = strconv.FormatInt(getServerTime(), 10)

	// Create the query string
	var queryParams []string
	for key, value := range params {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
	}
	queryString := strings.Join(queryParams, "&")
	
	// Construct the request URL
	
	// Sign the request using HMAC-SHA256
	signature := signRequest(queryString)

	url := baseURL + path+"?"+queryString+"&signature="+signature
	fmt.Println(url)
	// Create a new request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Add API key to request headers
	req.Header.Add("X-MBX-APIKEY", apiKey)


	// Send the request
	return client.Do(req)
}

func signRequest(queryString string) string {
	// Create a new HMAC-SHA256 hasher
	hasher := hmac.New(sha256.New, []byte(secretKey))

	// Write the query string to the hasher
	hasher.Write([]byte(queryString))

	// Get the resulting HMAC-SHA256 signature
	signature := hex.EncodeToString(hasher.Sum(nil))

	return signature
}

func getServerTime() int64 {

	headers := map[string][]string{
		"Accept": []string{"application/json"},
	}
	req, _ := http.NewRequest("GET", "https://api.binance.th/api/v1/time", nil)
	req.Header = headers
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	var respMap map[string]int64
	json.Unmarshal(body, &respMap)
	return respMap["serverTime"]
}
