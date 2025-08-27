
if [ -z "$3" ]; then
  echo "Usage: $0 <docker-image> <external-port> <internal-port>" >&2
  exit 1
fi

if ! docker image inspect $1 > /dev/null 2>&1; then
  echo "Building Docker image..." >&2
  docker build -t $1 . >&2
else
  echo "Docker image already exists. Skipping build." >&2
fi  

docker run -d -p $2:$3 $1