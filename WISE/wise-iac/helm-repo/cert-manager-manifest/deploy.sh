## reference link
#
#https://docs.bitnami.com/tutorials/secure-kubernetes-services-with-ingress-tls-letsencrypt/
#

#!/bin/bash

## check helm version is 3+ 
echo "confirm if helm version is v3.1+"
helm version

echo "confim if kubectl is installed"
kubectl version

##  pass cluster context to script as arg
# kubectl config use-context $1

## adding helm jetstack repo
helm repo add jetstack https://charts.jetstack.io

## creating cert-manager namespace
kubectl create namespace cert-manager

# install CRD for cert manager, webhook etc. 
## download link https://github.com/jetstack/cert-manager/releases/download/v0.14.1/cert-manager.crds.yaml
kubectl apply -f cert-manager.crds.yaml 

# install helm chart
echo "Installing helm chart ..... ... ... .. .. .. . . . . ."
helm install cert-manager --namespace cert-manager jetstack/cert-manager --version v0.14.1

## install letsencrypt issuer 

kubectl apply -f cluster-issuer-prod.yaml 
kubectl apply -f cluster-issuer-staging.yaml
