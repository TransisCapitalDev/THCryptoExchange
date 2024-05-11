import hmac
import hashlib
import requests
import json
import time

class BinanceThOrderObj:
    def __init__(self, symbol=None, orderId=None, clientOrderId=None, price=None,
                 origQty=None, executedQty=None, cumulativeQuoteQty=None,
                 status=None, timeInForce=None, type=None, side=None,
                 stopPrice=None, time=None, updateTime=None, isWorking=None,
                 origQuoteOrderQty=None):
        self.symbol = symbol
        self.orderId = orderId
        self.clientOrderId = clientOrderId
        self.price = price
        self.origQty = origQty
        self.executedQty = executedQty
        self.cumulativeQuoteQty = cumulativeQuoteQty
        self.status = status
        self.timeInForce = timeInForce
        self.type = type
        self.side = side
        self.stopPrice = stopPrice
        self.time = time
        self.updateTime = updateTime
        self.isWorking = isWorking
        self.origQuoteOrderQty = origQuoteOrderQty

def main():
    # Binance API endpoint
    endpoint = "https://api.binance.th/api/v1/allOrders"

    # Replace these with your actual API key and secret key
    api_key = ""
    secret_key = ""

    # Create a timestamp in milliseconds
    timestamp = get_server_time()

    dif = (time.time_ns() // 10**6) - timestamp

    # Create a query string with required parameters
    query_string = f"symbol=BTCTHB&timestamp={timestamp}"

    # Create a new HMAC-SHA256 hasher
    hasher = hmac.new(secret_key.encode(), query_string.encode(), hashlib.sha256)

    # Get the resulting HMAC-SHA256 signature
    signature = hasher.hexdigest()

    # Construct the request URL with the query string and signature
    url = f"{endpoint}?{query_string}&signature={signature}"

    # Create headers
    headers = {"X-MBX-APIKEY": api_key}

    # Send the request
    resp = requests.get(url, headers=headers)

    print(resp.status_code, resp.reason)
    print("=========")
    print(resp.text)


def get_server_time():
    url = "https://api.binance.th/api/v1/time"
    resp = requests.get(url)
    resp_json = resp.json()
    return resp_json["serverTime"]

if __name__ == "__main__":
    main()
