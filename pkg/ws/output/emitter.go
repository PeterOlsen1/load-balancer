package output

import (
	"encoding/json"
	"load-balancer/pkg/config"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"time"
)

type Emitter struct {
	LockedConn *types.LockedConnection
}

func getBaseEvent(eventType string) BaseEvent {
	return BaseEvent{
		Time: time.Now().Format(time.RFC3339),
		Type: eventType,
	}
}

func (s *Emitter) SendMessage(message string) error {
	if s.LockedConn == nil || !config.Config.Emitter.Enabled {
		return nil
	}

	s.LockedConn.Lock.Lock()
	defer s.LockedConn.Lock.Unlock()

	err := s.LockedConn.Conn.WriteMessage(1, []byte(message))
	if err != nil {
		logger.Err("Sending websocket message", err)
	}
	return err
}

func (s *Emitter) Request(conn *types.Connection) error {
	j, err := json.Marshal(Request{
		BaseEvent: getBaseEvent("request"),
		IP:        conn.Request.RemoteAddr,
		Method:    conn.Request.Method,
		Path:      conn.Request.URL.Path,
		UserAgent: conn.Request.UserAgent(),
	})

	if err != nil {
		logger.Err("WS: marshalling request json", err)
	}

	return s.SendMessage(string(j))
}

func (s *Emitter) Proxy(path string, proxiedTo string, ip string) error {
	j, err := json.Marshal(Proxy{
		BaseEvent: getBaseEvent("proxy"),
		IP:        ip,
		Path:      path,
		ProxiedTo: proxiedTo,
	})

	if err != nil {
		logger.Err("WS: marshalling proxy json", err)
	}

	return s.SendMessage(string(j))
}

func (s *Emitter) Health(status string, address string, respTime float32) error {
	j, err := json.Marshal(Health{
		BaseEvent:    getBaseEvent("health"),
		Status:       status,
		Address:      address,
		ResponseTime: respTime,
	})

	if err != nil {
		logger.Err("WS: marshalling health json", err)
	}

	return s.SendMessage(string(j))
}

func (s *Emitter) ContainerStart(containerID string) error {
	j, err := json.Marshal(ContainerStart{
		BaseEvent:   getBaseEvent("container_start"),
		ContainerID: containerID,
	})

	if err != nil {
		logger.Err("WS: marshalling container start json", err)
	}

	return s.SendMessage(string(j))
}

func (s *Emitter) ContainerStop(containerID string) error {
	j, err := json.Marshal(ContainerStop{
		BaseEvent:   getBaseEvent("container_stop"),
		ContainerID: containerID,
	})

	if err != nil {
		logger.Err("WS: marshalling container stop json", err)
	}

	return s.SendMessage(string(j))
}

func (s *Emitter) Info(message string) error {
	j, err := json.Marshal(Info{
		BaseEvent: getBaseEvent("info"),
		Message:   message,
	})

	if err != nil {
		logger.Err("WS: marshalling info json", err)
	}

	return s.SendMessage(string(j))
}

func (s *Emitter) Error(message string, err error) error {
	j, err := json.Marshal(Error{
		BaseEvent: getBaseEvent("error"),
		Message:   message,
		Error:     err,
	})

	if err != nil {
		logger.Err("WS: marshalling error json", err)
	}

	return s.SendMessage(string(j))
}
