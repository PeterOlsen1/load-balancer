# load-balancer
Trying out new things. Learning Go. Balancing loads. Scaling horizontally ⚖️

## current work
  
Testing:
* How do we make sure a node is unreachable before removing it
  * Close request queue, fulfill all waiting, close node
* Horizontal scaling
  * Update the health of a container as soon as it is good to go
  * Add container pool for faster scaling
* Test with multiple routes, `getRouteConfig` method

Working on:
* What to do when container queue is full?
  * Pop some off of the back
* How to stop non container nodes

Later issues:
* Web frontend
* Fix reciever methods in balancer

Future ideas:
* Batch process requests from node queue?
* Add rate tracking?
* Deploying to AWS or something?

## architecture planning

Questions:
* When to spin up a new container
* What metrics should we send to frontend?

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

### items completed

* Fix memory leak with idle looping in queue watching. don't use default case!
* Add a request queue for each node
* Config
  * use docker SDK
  * More config rules?
    * URLs, rules for proxy, etc.
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

iterating on this:
```YAML
routes:
  - path: /api/*
    docker_image: api-server
    strategy: least_connections
    servers:
    - url: http://aws.something.com
    - url: http://hello.com
    - ...
    
```
