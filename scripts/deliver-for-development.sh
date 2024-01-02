#!/bin/bash

echo 'The following command lines builds your Golang application for'
echo 'development environment'

set -x

docker build -t gcr.io/bexs-platform/temis-compliance:${SHORT_COMMIT} ../.
docker push gcr.io/bexs-platform/temis-compliance:${SHORT_COMMIT}
gcloud config set project bexs-digital-dev
gcloud secrets versions access latest --secret=temis --format='get(payload.data)' | tr '_-' '/+' | base64 -d > ./secrets.encrypted
gcloud kms decrypt --location 'global' --keyring 'Digital' --key 'temis' --ciphertext-file ./secrets.encrypted --plaintext-file ../.kubernetes/develop/secrets.yaml
gcloud config set compute/zone southamerica-east1
gcloud config set compute/region southamerica-east1
gcloud config set project bexs-develop
gcloud container clusters get-credentials develop
export https_proxy=10.40.15.195:80
cp -r ../.kubernetes/yaml/* ../.kubernetes/develop
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-develop:southamerica-east1:temis/g" ../.kubernetes/develop/deploy-temis-compliance-api-ext.yaml
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-develop:southamerica-east1:temis/g" ../.kubernetes/develop/deploy-temis-compliance-api-int.yaml
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-develop:southamerica-east1:temis/g" ../.kubernetes/develop/deploy-temis-compliance-async.yaml
sed -i -e "s/ENV/dev/g" ../.kubernetes/develop/ingress-compliance-api.yaml
sed -i -e "s/ENV/dev/g" ../.kubernetes/develop/vs-compliance-external.yaml

sed -i -e "s/SERVICE_ACCOUNT_NAME/temis-compliance/g" \
       -e "s/SERVICE_ACCOUNT_PROJECT/bexs-develop/g" ../.kubernetes/develop/serviceAccount.yaml

kubectl apply -f ../.kubernetes/develop/

set +x
