
if [ -z "$1" ]; then
  echo "Usage: $0 <port>" >&2
  exit 1
fi

if ! docker image inspect node-server > /dev/null 2>&1; then
  echo "Building Docker image..." >&2
  docker build -t node-server . >&2
else
  echo "Docker image already exists. Skipping build." >&2
fi  

docker run -d -p $1:3000 node-server