package ws

//emit events to the frontend
import (
	"load-balancer/pkg/logger"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Err("Failed to upgrade connection", err)
		return
	}
	defer conn.Close()

	for {
		messageType, body, err := conn.ReadMessage()
		if err != nil {
			logger.Err("Reading from websocket", err)
			return
		}
		logger.WsRequest(body)

		if err := conn.WriteMessage(messageType, body); err != nil {
			logger.Err("websocket write", err)
			return
		}
	}
}
