package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type BinanceThOrderObj struct {
	Symbol             string `json:"symbol,omitempty"`
	OrderId            int    `json:"orderId,omitempty"`
	ClientOrderId      string `json:"clientOrderId,omitempty"`
	Price              string `json:"price,omitempty"`
	OrigQty            string `json:"origQty,omitempty"`
	ExecutedQty        string `json:"executedQty,omitempty"`
	CumulativeQuoteQty string `json:"cumulativeQuoteQty,omitempty"`
	Status             string `json:"status,omitempty"`
	TimeInForce        string `json:"timeInForce,omitempty"`
	Type               string `json:"type,omitempty"`
	Side               string `json:"side,omitempty"`
	StopPrice          string `json:"stopPrice,omitempty"`
	Time               int64  `json:"time,omitempty"`
	UpdateTime         int64  `json:"updateTime,omitempty"`
	IsWorking          bool   `json:"isWorking,omitempty"`
	OrigQuoteOrderQty  string `json:"origQuoteOrderQty,omitempty"`
}

func main() {
	// Binance API endpoint
	endpoint := "https://api.binance.th/api/v1/allOrders"

	// Replace these with your actual API key and secret key
	apiKey := " "
	secretKey := " "
	// Create a timestamp in milliseconds
	timestamp := getServerTime()
	fmt.Println((time.Now().UnixNano() / int64(time.Millisecond)))

	dif := (time.Now().UnixNano() / int64(time.Millisecond)) - timestamp
	fmt.Println(dif)

	// Create a query string with required parameters
	queryString := fmt.Sprintf("symbol=BTCTHB&timestamp=%d", timestamp)

	// Create a new HMAC-SHA256 hasher
	hasher := hmac.New(sha256.New, []byte(secretKey))

	// Write the query string bytes to the hasher
	hasher.Write([]byte(queryString))

	// Get the resulting HMAC-SHA256 signature
	signature := hex.EncodeToString(hasher.Sum(nil))

	// Construct the request URL with the query string and signature
	url := fmt.Sprintf("%s?%s&signature=%s", endpoint, queryString, signature)

	// Create a new HTTP client
	client := http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the API key header
	req.Header.Add("X-MBX-APIKEY", apiKey)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, resp.Status)
	fmt.Println("=========")
	fmt.Println(string(body))
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	var respMap []BinanceThOrderObj
	json.Unmarshal(body, &respMap)
	// Print the response body
	fmt.Println(respMap)
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
