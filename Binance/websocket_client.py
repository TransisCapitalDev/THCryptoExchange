import asyncio
import websockets
import signal
import json
import time
import logging

logging.basicConfig(level=logging.INFO)

async def process_msg(message):
    logging.info(f"Received: {message}")

async def main():
    url = "wss://stream-th.2meta.app/stream"
    logging.info(f"Connecting to {url}")

    async with websockets.connect(url) as websocket:
        done = asyncio.Event()

        subscribe_message = {
            "method": "SUBSCRIBE",
            "params": [
                "btcthb@depth20",
            ],
            "id": 1,
        }

        # Send the subscription message
        await websocket.send(json.dumps(subscribe_message))

        async def read_messages():
            try:
                async for message in websocket:
                    await process_msg(message)
            except websockets.ConnectionClosed:
                logging.info("Connection closed, exiting reader loop")
                done.set()

        async def send_ticker():
            try:
                while not done.is_set():
                    timestamp = str(time.time())
                    await websocket.send(timestamp)
                    await asyncio.sleep(1)
            except websockets.ConnectionClosed:
                logging.info("Connection closed, exiting sender loop")
                done.set()

        # Handle interrupts
        loop = asyncio.get_event_loop()
        stop = asyncio.Future()

        def handle_interrupt():
            if not stop.done():
                stop.set_result(None)

        loop.add_signal_handler(signal.SIGINT, handle_interrupt)
        loop.add_signal_handler(signal.SIGTERM, handle_interrupt)

        await asyncio.gather(read_messages(), send_ticker(), stop)

        logging.info("Interrupt received, closing connection")
        await websocket.close()

if __name__ == "__main__":
    asyncio.run(main())
