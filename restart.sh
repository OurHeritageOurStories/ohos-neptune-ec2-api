if [ $# -eq 0 ]
  then
    echo "No arguments supplied. You must supply the environment you want to run; either 'live' or 'test'"
    exit 1
fi
export CURRENT_ENV=$1
docker compose down
docker rmi -f $(docker images -aq)
docker compose up -d
sleep 20
curl --request POST \
  --url http://localhost:8001/config \
  --header 'Content-Type: text/yaml' \
  --data '_format_version: "2.1"
_transform: true

services: 

  - name: goapi
    url: http://goapi:9000/
    routes:
      - name: goapi
        paths:
          - /api
        preserve_host: true
        strip_path: true
        methods:
          - GET
          - POST
  - name: frontend
    url: http://frontend:3000/
    routes:
      - name: frontend
        paths:
          - /
        preserve_host: true
        strip_path: true
        methods:
          - GET
  '
if [[ $1 = "test" ]] 
  then
    cd go-api/
    go test -v
    cd ../
else 
  echo "Not running tests"
fi
