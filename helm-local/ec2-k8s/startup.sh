cd go
docker build -t goapi .
cd ../

minikube start

helm install kong kong --set service.port=5000
helm install goapi goapi

sleep 60

kubectl port-forward svc/kong :5000 --address='0.0.0.0'