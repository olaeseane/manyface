#!/bin/sh


kubectl delete -f conduit-configmap.yaml
kubectl delete -f conduit-deployment.yaml
sleep 2
kubectl apply -f conduit-configmap.yaml
kubectl apply -f conduit-deployment.yaml


# kubectl apply -f conduit-configmap.yaml
# kubectl apply -f conduit-pvc.yaml
# kubectl apply -f conduit-deployment.yaml
# kubectl apply -f conduit-clusterip.yaml

# sleep 10

# kubectl apply -f manyface-configmap.yaml
# kubectl apply -f manyface-pvc.yaml
# kubectl apply -f manyface-deployment.yaml
# kubectl apply -f manyface-clusterip.yaml

# kubectl apply -f test-pod.yaml
