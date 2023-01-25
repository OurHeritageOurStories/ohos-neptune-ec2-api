start minikube

#helm install miiify miiify
helm install blazegraph blazegraph
helm install go-api go-api
helm install wqds wqds
helm install react-frontend react-frontend
#helm install kong kong --set service.port=5005
#helm install iiif-generator iiif-generator

minikube tunnel
