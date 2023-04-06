docker compose down
export CURRENT_ENV="test"
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
cd go-api/
go test -v
cd ../

# TODO if tests work, do the below, else do something else
docker compose down
export CURRENT_ENV="live"
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