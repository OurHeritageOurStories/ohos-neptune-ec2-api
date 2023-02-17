minikube start

helm install kong kong --set service.port=5000
helm install miiify miiify
helm install blazegraph blazegraph
#helm install http-echo http-echo
helm install echotest echotest

sleep 60

kubectl port-forward svc/kong :5000 --address='0.0.0.0'