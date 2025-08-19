# load-balancer
Trying out new things. Learning Go. Balancing loads ⚖️

## architecture planning

Questions:
* When to spin up a new container
* What metrics should we send to frontend?

Working on:
* Config
  * Add different route options? for exmaple:
 ```YAML
routes:
  - path: /api/*
    backend: api_servers
  - path: /static/*
    backend: static_servers

backends:
  api_servers:
    strategy: least_connections
    servers:
      - url: http://10.0.0.1:8080
      - url: http://10.0.0.2:8080
```
  * More config rules?
    * URLs, rules for proxy, etc.
  * Replace hardcoded values with config
* More research on when to update running containers
* Load URLs into nodes from CSV file when starting, or at least keep them for later use
* Test out websocket connection more

Future ideas:
* Add rate tracking?
* Testing framework (really load it)

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
    * least connections (+ weighted)
    * compute based ?
    * IP hash

Things to research:
* Good load balancing algorithms / strategies
  * At what point do we decide to spin up a new instance
  * How do we decide which instance to send requests to?

### frontend thoughts
* Communicate to the backend with websockets
  * Backend shoots logs to the frontend
  * Logs + metrics about responses
* Allow for user to manually start / update nodes
