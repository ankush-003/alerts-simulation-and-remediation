#!/bin/bash

kubectl apply -f secrets.yaml
kubectl apply -f rest-deployment.yaml
kubectl apply -f rest-server-service.yaml
kubectl apply -f rules-engine-deployment.yaml
kubectl apply -f sim-deployment.yaml

# Check the status of the deployments
kubectl get deployments
kubectl get services
kubectl get pods

# Check the minikube dashboard
minikube dashboard