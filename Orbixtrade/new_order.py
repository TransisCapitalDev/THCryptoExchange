import hmac
import hashlib
import json
import time
import requests
from collections import OrderedDict


api_key = " "
secret_key = " "
baseURL = "https://www.orbixtrade.com/api/"
orders = "orders/"

class Option:
    def __init__(self, type_, pair, side, nonce, price=None, amount=None):
        self.type = type_
        self.pair = pair
        self.side = side
        self.price = price
        self.amount = amount
        self.nonce = nonce

def main():
    now = time.time()
    unix_timestamp = int(now)
    options = Option(
        type_="limit",
        pair="btc_thb",
        side="buy",
        nonce=unix_timestamp,
        price="2000000",
        amount="0.0002"
    )

    params = options.__dict__
    bb = json.dumps(params)
    qs = construct_query_string_with_prefix(params, "")
    print("payload", qs)
    sig = sign_request(qs)

    headers = {
        "Authorization": f"TDAX-API {api_key}",
        "Signature": sig,
        "Content-Type": "application/json"
    }

    response = requests.post(f"{baseURL}{orders}", data=bb, headers=headers)
    print(response)

def sign_request(query_string):
    mac = hmac.new(secret_key.encode(), query_string.encode(), hashlib.sha512)
    sig = mac.hexdigest()
    return sig

def construct_query_string_with_prefix(params, prefix):
    keys = sorted(params.keys())
    qs = ""

    for i, k in enumerate(keys):
        v = params[k]

        if isinstance(v, dict):
            qs += construct_query_string_with_prefix(v, k)
        elif isinstance(v, list):
            nested_map = {str(i): val for i, val in enumerate(v)}
            qs += construct_query_string_with_prefix(nested_map, k)
        else:
            if prefix:
                if isinstance(v, float):
                    qs += f"{prefix}[{k}]={v:.0f}"
                else:
                    qs += f"{prefix}[{k}]={v}"
            else:
                if isinstance(v, float):
                    qs += f"{k}={v:.0f}"
                else:
                    qs += f"{k}={v}"

        if i != len(keys) - 1:
            qs += "&"

    return qs

if __name__ == "__main__":
    main()
