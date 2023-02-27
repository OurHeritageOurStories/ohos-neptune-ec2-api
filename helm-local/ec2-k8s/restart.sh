docker compose down
docker rmi -f $(docker images -aq)
docker compose up -d
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
          - /
        preserve_host: true
        strip_path: true
        methods:
          - GET
          - POST
  '
  