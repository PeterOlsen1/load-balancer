# load-balancer
Trying out new things. Learning Go. Balancing loads 🤑

### architecture planning

Basic idea: we want to have some sort of web service, where requests are recieved, sent to a shared queue, and then handled accordingly

We need loops (within goroutines) for:
* adding to shared queue
* removing from queue and distributing to servers

Goroutines for:
* Checking node health
* Spinning up new docker container instances

Things to research:
* Good load balancing algorithms / strategies
  * At what point do we decide to spin up a new instance
  * How do we decide which instance to send requests to?
* How the hell to use docker

Other stuff:
* Make a logger (good practice)
* Decide what sort of stuff to serve on this web server
  * Simple JS thing?
* Websocket connection to frontend for live monitoring?
  * Tauri app?