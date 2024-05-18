package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Option struct {
	Type   string `json:"type"` // possible value ['limit','market']
	Pair   string `json:"pair"`
	Side   string `json:"side"` // possible value ['buy','sell']
	Price  string `json:"price,omitempty"`
	Amount string `json:"amount,omitempty"`
	Nonce  int64  `json:"nonce"`
}

var (
	apiKey    = " "
	secretKey = " "
	baseURL   = "https://www.orbixtrade.com/api/"
	orders    = "orders/"
)

func main() {

	now := time.Now()
	unixTimestame := now.Unix()
	options := Option{
		Type:   "limit",
		Pair:   "btc_thb",
		Side:   "buy",
		Nonce:  unixTimestame,
		Price:  "2000000",
		Amount: "0.0002",
	}

	params := map[string]interface{}{}
	bb, _ := json.Marshal(options)
	json.Unmarshal(bb, &params)
	str := strings.NewReader(string(bb))

	qs := constructQueryStringWithPrefix(params, "")
	fmt.Println("payload", qs)
	sig := signRequest(qs)

	req, err := http.NewRequest("POST", baseURL+orders, str)
	if err != nil {
		fmt.Println(err.Error())

	}
	// Add API key to request headers
	req.Header.Add("Authorization", "TDAX-API "+apiKey)
	req.Header.Add("Signature", sig)
	req.Header.Add("Content-Type", "application/json")

	clientHttp := http.Client{}
	resp, _ := clientHttp.Do(req)
	fmt.Println(resp)

}

func signRequest(queryString string) string {
	mac := hmac.New(sha512.New, []byte(secretKey))
	mac.Write([]byte(queryString))
	sig := mac.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(sig)))
	hex.Encode(dst, sig)
	return string(dst)
}

func constructQueryStringWithPrefix(params map[string]interface{}, prefix string) string {
	var keys []string

	for k, _ := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var qs string

	for i, k := range keys {
		v := params[k]

		if nestedParams, ok := v.(map[string]interface{}); ok {
			qs += constructQueryStringWithPrefix(nestedParams, k)
		} else if array, ok := v.([]interface{}); ok {
			logrus.Errorf("bar")
			nestedMap := map[string]interface{}{}
			for i, v := range array {
				nestedMap[fmt.Sprintf("%d", i)] = v
			}
			qs += constructQueryStringWithPrefix(nestedMap, k)
		} else {
			if prefix == "" {
				if _, ok := v.(float64); ok {
					qs += fmt.Sprintf("%s=%.f", k, v)
				} else {
					qs += fmt.Sprintf("%s=%v", k, v)
				}
			} else {
				if _, ok := v.(float64); ok {
					qs += fmt.Sprintf("%s[%s]=%.f", prefix, k, v)
				} else {
					qs += fmt.Sprintf("%s[%s]=%v", prefix, k, v)
				}
			}
		}

		if i != len(keys)-1 {
			qs += "&"
		}
	}

	return qs
}
