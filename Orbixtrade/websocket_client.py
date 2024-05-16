import asyncio
import websockets
import signal
import json
import time

async def process_msg(message):
    print(f"Received: {message}")

async def main():
    url = "wss://www.orbixtrade.com/ws/stream?streams=!miniTicker@arr@3000ms"
    print(f"Connecting to {url}")

    async with websockets.connect(url) as websocket:
        subscribe_message = {
            "method": "SUBSCRIBE",
            "params": [
                # "btcthb@trade",
                "btc_thb@depth",
                "btc_thb@aggTrade",
            ],
            "id": 2,
        }

        # Send the subscription message
        await websocket.send(json.dumps(subscribe_message))

        async def read_messages():
            async for message in websocket:
                await process_msg(message)

        # Handle interrupts
        loop = asyncio.get_event_loop()
        stop = asyncio.Future()

        def handle_interrupt():
            if not stop.done():
                stop.set_result(None)

        loop.add_signal_handler(signal.SIGINT, handle_interrupt)

        # Send a message every second
        async def send_ticker():
            while not stop.done():
                await websocket.send(json.dumps({"type": "ping", "time": time.time()}))
                await asyncio.sleep(1)

        await asyncio.gather(read_messages(), send_ticker(), stop)

        # Cleanly close the connection
        close_message = json.dumps({"type": "close", "time": time.time()})
        await websocket.send(close_message)

if __name__ == "__main__":
    asyncio.run(main())
