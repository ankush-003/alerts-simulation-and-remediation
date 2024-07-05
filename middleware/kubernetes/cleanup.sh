#!/bin/bash

# Delete the deployments and services defined in the YAML files
kubectl delete -f rest-server-setup.yaml
kubectl delete -f rules-engine-setup.yaml
kubectl delete -f sim-setup.yaml
kubectl delete -f dashboard-setup.yaml

# Delete the secrets
kubectl delete -f secrets.yaml

# Optionally, delete specific resources by name if they are not covered in the above files
kubectl delete deployment simulator
# kubectl delete service <service-name>

# Check the status of the deployments, services, and pods to ensure they are deleted
kubectl get deployments
kubectl get services
kubectl get pods

# Shutdown Minikube
# minikube stop
