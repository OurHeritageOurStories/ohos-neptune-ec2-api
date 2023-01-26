minikube start

helm install kong kong --set service.port=5000
helm install miiify miiify
helm install blazegraph blazegraph

sleep 60

kubectl port-forward svc/kong 5000:5000