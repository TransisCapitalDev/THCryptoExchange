import hmac
import hashlib
import requests
import json
import time


api_key = " "
secret_key = " "
base_url = "https://api.binance.th"

def main():
    # Set the request parameters
    params = {
        "symbol": "BTCTHB",
        "side": "BUY",
        "type": "LIMIT",
        "timeInForce": "GTC",
        "quantity": "0.0002",
        "price": "2000000",
    }

    # Send the request
    resp = send_request("POST", "/api/v1/order", params)
    if resp.status_code != 200:
        print("Error sending request:", resp.text)
        return

    # Print the response body
    print(resp.text)

def send_request(method, path, params):
    # Add timestamp to request params
    params["timestamp"] = str(get_server_time())

    # Construct the query string
    query_string = "&".join([f"{key}={value}" for key, value in params.items()])

    # Sign the request using HMAC-SHA256
    signature = sign_request(query_string)

    url = f"{base_url}{path}?{query_string}&signature={signature}"
    print(url)

    # Create headers
    headers = {"X-MBX-APIKEY": api_key}

    # Send the request
    if method == "GET":
        return requests.get(url, headers=headers)
    elif method == "POST":
        return requests.post(url, headers=headers)

def sign_request(query_string):
    # Create a new HMAC-SHA256 hasher
    hasher = hmac.new(secret_key.encode(), query_string.encode(), hashlib.sha256)

    # Get the resulting HMAC-SHA256 signature
    return hasher.hexdigest()

def get_server_time():
    url = f"{base_url}/api/v1/time"
    resp = requests.get(url)
    resp_json = resp.json()
    return resp_json["serverTime"]

if __name__ == "__main__":
    main()
