package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	url := "wss://www.orbixtrade.com/ws/stream?streams=!miniTicker@arr@3000ms"
	fmt.Printf("Connecting to %s\n", url)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	subscribeMessage := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{
			//"btcthb@trade",
			"btc_thb@depth",
			"btc_thb@aggTrade",
		},
		"id": 2,
	}

	// Send the subscription message
	err = c.WriteJSON(subscribeMessage)
	if err != nil {
		log.Fatalf("Error sending subscription message: %v", err)
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		// To do process message
		ProcessMsg(message)

	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("Write error:", err)
				return
			}
		case <-interrupt:
			log.Println("Interrupt received, closing connection")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
func ProcessMsg(message []byte) {
	log.Printf("Received: %s", message)
}
