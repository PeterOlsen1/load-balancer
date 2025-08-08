if [ -z "$1" ]; then
  echo "Usage: $0 <port>"
  exit 1
fi

if ! docker image inspect node-server > /dev/null 2>&1; then
  echo "Building Docker image..."
  docker build -t node-server .
else
  echo "Docker image already exists. Skipping build."
fi  

echo "Running server on port $1"
docker run -p $1:3000 node-server