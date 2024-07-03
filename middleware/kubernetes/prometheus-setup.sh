#!/bin/env bash
# Reference: https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
# make sure to have helm installed

# Add the prometheus-community and stable helm repos
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add stable https://charts.helm.sh/stable
helm repo update

# Install chart
helm install prometheus prometheus-community/kube-prometheus-stack
# link to chart https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
