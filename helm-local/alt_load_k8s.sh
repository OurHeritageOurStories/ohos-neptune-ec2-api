minikube start

helm install kong kong --set service.port=5000
helm install miiify miiify

sleep 300

kubectl port-forward svc/kong 5000:5000