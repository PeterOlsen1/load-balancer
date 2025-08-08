if [ -z "$1" ]; then
  echo "Usage: $0 <port>"
  exit 1
fi

docker build -t node-server .  

echo "Running server on port $1"

docker run -p $1:3000 node-server