package ws

//emit events to the frontend
import (
	"load-balancer/pkg/types"
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
var EventReciever input.Receiver = input.InitReceiver()

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Err("Failed to upgrade connection", err)
		return
	}
	defer conn.Close()

	lockedConn := types.LockedConnection{
		Conn: conn,
	}

	logger.WsConnect(r)
	EventEmitter = output.Emitter{
		LockedConn: &lockedConn,
	}
	defer func() {
		//remove connection at the end
		EventEmitter.LockedConn = nil
	}()

	for {
		_, body, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) ||
				websocket.IsUnexpectedCloseError(err) {
				logger.WsClose(r)
			} else {
				logger.Err("Reading from websocket", err)
			}
			return
		}
		logger.WsRequest(body, r.RemoteAddr)

		bytes, err := EventReciever.HandleWsRequest(body)
		if err != nil {
			//errror should be logged within function
			EventEmitter.Error("Constructing response", err)
			conn.WriteMessage(0, []byte("{}"))
			continue
		}

		lockedConn.Lock.Lock()
		defer lockedConn.Lock.Unlock()
		err = conn.WriteMessage(1, bytes)
		if err != nil {
			logger.Err("Writing websocket response", err)
		}
	}
}
