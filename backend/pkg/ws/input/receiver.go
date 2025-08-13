package input

import (
	"encoding/json"
	"fmt"
	"load-balancer/pkg/logger"
	"os"
)

type ReceiverFunction func(body string) (string, error)

type Receiver struct {
	eventMap map[string]ReceiverFunction
}

func (r *Receiver) AddEventHandler(eventType string, f ReceiverFunction) {
	existingFunc := r.eventMap[eventType]
	if existingFunc != nil {
		fmt.Println("Attempt to add two event handlers!")
		os.Exit(1)
		return
	}

	r.eventMap[eventType] = f
}

func (r *Receiver) HandleWsRequest(body []byte) ([]byte, error) {
	//call hooks set in place by receiver files
	var base BaseEvent

	err := json.Unmarshal(body, &base)
	if err != nil {
		logger.Err("Unmarshaling ws request", err)
		return nil, err
	}

	messageType := base.Type
	receiverFunc := r.eventMap[messageType]
	if receiverFunc == nil {
		err := fmt.Errorf("no receiver hook set for ws input type: %s", messageType)
		logger.Err("Handling ws request", err)
		return nil, err
	}

	bytes, err := receiverFunc(string(body))
	return []byte(bytes), err
}
