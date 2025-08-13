package ws

//emit events to the frontend
import (
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws/input"
	"load-balancer/pkg/ws/output"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var EventEmitter output.Emitter
var EventReciever input.Receiver

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Err("Failed to upgrade connection", err)
		return
	}
	defer conn.Close()

	EventEmitter = output.Emitter{
		Conn: conn,
	}

	for {
		_, body, err := conn.ReadMessage()
		if err != nil {
			logger.Err("Reading from websocket", err)
			continue
		}
		logger.WsRequest(body)

		bytes, err := EventReciever.HandleWsRequest(body)
		if err != nil {
			//errror should be logged within function
			EventEmitter.Error("Constructing response", err)
			conn.WriteMessage(0, []byte("{}"))
			continue
		}

		err = conn.WriteMessage(1, bytes)
		if err != nil {
			logger.Err("Writing websocket response", err)
		}
	}
}
