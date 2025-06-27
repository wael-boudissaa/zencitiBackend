package user

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (h *Handler) ClientLocationWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	for {
		var msg struct {
			IdClient  string  `json:"idClient"`
			Longitude float64 `json:"longitude"`
			Latitude  float64 `json:"latitude"`
		}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		h.store.UpdateClientLocation(msg.IdClient, msg.Longitude, msg.Latitude)
	}
}
