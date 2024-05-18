import hmac
import hashlib
import json
import requests
from collections import OrderedDict


api_key = " "
secret_key = " "
base_url = "https://www.orbixtrade.com/api/"
orders = "orders/user"

class Option:
    def __init__(self, pair, side, limit, offset, status=""):
        self.status = status
        self.pair = pair
        self.side = side
        self.limit = limit
        self.offset = offset

def main():
    options = Option(
        pair="btc_thb",
        side="buy",
        limit=20,
        offset=0,
    )

    params = options.__dict__
    bb = json.dumps(params)
    params = json.loads(bb)

    qs = construct_query_string_with_prefix(params, "")
    sig = sign_request("")
    print("payload ", f"{base_url}{orders}?{qs}")
    headers = {
        "Authorization": f"TDAX-API {api_key}",
        "Signature": sig,
        "Content-Type": "application/json"
    }

    response = requests.get(f"{base_url}{orders}?{qs}", headers=headers)
    print(response.status_code)
    print(response.text)

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
