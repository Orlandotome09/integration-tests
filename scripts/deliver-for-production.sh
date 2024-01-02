#!/bin/bash

echo 'The following command lines builds your Golang application for'
echo 'production environment'

set -x

docker build -t gcr.io/bexs-platform/temis-compliance:${SHORT_COMMIT} ../.
docker push gcr.io/bexs-platform/temis-compliance:${SHORT_COMMIT}
gcloud config set project bexs-digital
gcloud secrets versions access latest --secret=temis --format='get(payload.data)' | tr '_-' '/+' | base64 -d > ./secrets.encrypted
gcloud kms decrypt --location 'global' --keyring 'Digital' --key 'temis' --ciphertext-file ./secrets.encrypted --plaintext-file ../.kubernetes/production/secrets.yaml
gcloud config set compute/zone southamerica-east1
gcloud config set compute/region southamerica-east1
gcloud config set project bexs-production
gcloud container clusters get-credentials production-shared
export https_proxy=10.60.0.28:80
cp -r ../.kubernetes/yaml/* ../.kubernetes/production
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-production:southamerica-east1:temis-prd/g" ../.kubernetes/production/deploy-temis-compliance-api-ext.yaml
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-production:southamerica-east1:temis-prd/g" ../.kubernetes/production/deploy-temis-compliance-api-int.yaml
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-production:southamerica-east1:temis-prd/g" ../.kubernetes/production/deploy-temis-compliance-async.yaml
sed -i -e "s/ENV/prd/g" ../.kubernetes/production/ingress-compliance-api.yaml
sed -i -e "s/ENV/prd/g" ../.kubernetes/production/vs-compliance-external.yaml

sed -i -e "s/SERVICE_ACCOUNT_NAME/temis-compliance/g" \
       -e "s/SERVICE_ACCOUNT_PROJECT/bexs-production/g" ../.kubernetes/production/serviceAccount.yaml

kubectl apply -f ../.kubernetes/production/

set +x
