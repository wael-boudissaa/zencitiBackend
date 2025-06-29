package restaurant

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
)

type Client struct {
	conn         *websocket.Conn
	restaurantID string
	timeSlot     string
	tableID      string
	send         chan []byte
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

var clients = make(map[string][]*Client)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	var statusTables types.GetStatusTables
	result := utils.ParseJson(r, &statusTables)
	if result != nil {
		utils.WriteError(w, http.StatusBadRequest, result)
		return
	}
	restaurantID := statusTables.RestaurantId
	timeSlot := statusTables.TimeSlot

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade:", err)
		return
	}

	client := &Client{conn: conn, restaurantID: restaurantID, timeSlot: timeSlot, send: make(chan []byte)}

	go readPump(client)
	go writePump(client)
}

func readPump(client *Client) {
	defer func() {
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break // Exit loop on error or disconnect
		}

		// You can process the message here if clients send any (optional)
		log.Printf("Received message: %s", message)
	}
}

func writePump(client *Client) {
	defer func() {
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write error:", err)
				return
			}
		}
	}
}
