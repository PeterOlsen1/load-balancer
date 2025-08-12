# load-balancer
Trying out new things. Learning Go. Balancing loads 🤑

## architecture planning

Basic idea: we want to have some sort of web service, where requests are recieved, sent to a shared queue, and then handled accordingly

Questions:
* How often do we run the /health check
* When to spin up a new container
* What metrics should we collect?

Working on:
* Occasionally poll nodes for health, do in separate goroutine

Future ideas:
* Websocket connection to frontend for live monitoring?
  * Tauri app?
* Configure this to allow users to upload their own docker images?
* Add rate tracking?
* Re-implement shared queue

### general thoughts
  
We need loops (within goroutines) for:
* adding to shared queue
* removing from queue and distributing to servers

Goroutines for:
* Checking node health
* Spinning up new docker container instances

Proxy process:
* Dequeue a connection
* Pick a node to send it to (configure this to allow for multiple balancing strategies)
  * This requires the balancer internal data structure to be modified
    * Hash table for 'address: { node, connectionData }' pairs
  * Balancing options
    * Round robin
    * least connections (weighted)
    * compute based
    * IP hash

Things to research:
* Good load balancing algorithms / strategies
  * At what point do we decide to spin up a new instance
  * How do we decide which instance to send requests to?