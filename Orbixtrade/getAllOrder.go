package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/sirupsen/logrus"
)

type Option struct {
	Status string `json:"status"`
	Pair   string `json:"pair"`
	Side   string `json:"side,omitempty"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

var (
	apiKey    = " "
	secretKey = " "
	baseURL   = "https://www.orbixtrade.com/api/"
	orders    = "orders/user"
)

func main() {

	options := Option{
		Pair:   "btc_thb",
		Side:   "buy",
		Limit:  20,
		Offset: 0,
		Status: "",
	}
	params := map[string]interface{}{}
	bb, _ := json.Marshal(options)
	json.Unmarshal(bb, &params)

	qs := constructQueryStringWithPrefix(params, "")
	sig := signRequest("")
	fmt.Println("payload ", baseURL+orders+"?"+qs)
	req, err := http.NewRequest("GET", baseURL+orders+"?"+qs, nil)
	if err != nil {
		fmt.Println(err.Error())

	}
	// Add API key to request headers
	req.Header.Add("Authorization", "TDAX-API "+apiKey)
	req.Header.Add("Signature", sig)
	req.Header.Add("Content-Type", "application/json")

	clientHttp := http.Client{}
	resp, _ := clientHttp.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())

	}
	fmt.Println(resp.Status)
	fmt.Println(string(body))

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
