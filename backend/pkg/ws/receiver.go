package ws

type ReceiverFunction func(body string, err error) (string, error)

func handleWsRequest(body string, err error) {
	//call hooks set in place by receiver files
}
