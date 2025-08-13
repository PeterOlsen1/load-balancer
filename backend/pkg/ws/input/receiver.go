package input

import (
	"fmt"
	"os"
)

type ReceiverFunction func(body string, err error) (string, error)

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

func (r *Receiver) HandleWsRequest(body string, err error) {
	//call hooks set in place by receiver files
}
