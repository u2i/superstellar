#!/bin/sh
gcloud auth activate-service-account --key-file /deployment_volume/service_account.json
gcloud config set project kubernetes-playground-195112
gcloud container clusters get-credentials u2i-cluster --zone=northamerica-northeast1-b
kubectl --namespace=superstellar set image deployment/superstellar-backend-deployment superstellar-backend=gcr.io/kubernetes-playground-195112/superstellar-backend:"$1"
kubectl --namespace=superstellar set image deployment/superstellar-frontend-deployment superstellar-frontend=gcr.io/kubernetes-playground-195112/superstellar-frontend:"$1"
