#!/bin/bash

minikube addons enable ingress

kubectl apply -f secrets.yaml
kubectl apply -f rest-server-setup.yaml
kubectl apply -f rules-engine-setup.yaml
kubectl apply -f sim-setup.yaml
kubectl apply -f dashboard-setup.yaml

# Check the status of the deployments
kubectl get deployments
kubectl get services
kubectl get pods

# Check the minikube dashboard
minikube dashboard