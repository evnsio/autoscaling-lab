#!/bin/bash

#
# Install stuff
#
minikube start --kubernetes-version v1.12.10
minikube addons enable metrics-server

helm install prometheus ./prometheus
helm install grafana ./grafana

kubectl apply -f vpa
./vpa/gencerts.sh

kubectl apply -f test_workloads