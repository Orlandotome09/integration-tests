#!/bin/bash

echo 'The following command lines builds your Golang application for'
echo 'sandbox environment!'

set -x

docker build -t gcr.io/bexs-platform/temis-compliance:${SHORT_COMMIT} ../.
docker push gcr.io/bexs-platform/temis-compliance:${SHORT_COMMIT}
gcloud config set compute/zone southamerica-east1
gcloud config set compute/region southamerica-east1
gcloud config set project bexs-sandbox
gcloud container clusters get-credentials sandbox
gcloud secrets versions access latest --secret=temis --format='get(payload.data)' | tr '_-' '/+' | base64 -d > ./secrets.encrypted
gcloud kms decrypt --location 'global' --keyring 'digital' --key 'temis' --ciphertext-file ./secrets.encrypted --plaintext-file ../.kubernetes/sandbox/secrets.yaml
cp -r ../.kubernetes/yaml/* ../.kubernetes/sandbox
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-sandbox:southamerica-east1:temis-sandbox/g" ../.kubernetes/sandbox/deploy-temis-compliance-api-ext.yaml
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-sandbox:southamerica-east1:temis-sandbox/g" ../.kubernetes/sandbox/deploy-temis-compliance-api-int.yaml
sed -i -e "s/SHORT_COMMIT/${SHORT_COMMIT}/g" -e "s/DB_NAME/bexs-sandbox:southamerica-east1:temis-sandbox/g" ../.kubernetes/sandbox/deploy-temis-compliance-async.yaml
sed -i -e "s/ENV/sandbox/g" ../.kubernetes/sandbox/ingress-compliance-api.yaml
sed -i -e "s/ENV/sandbox/g" ../.kubernetes/sandbox/vs-compliance-external.yaml

sed -i -e "s/SERVICE_ACCOUNT_NAME/temis-compliance/g" \
       -e "s/SERVICE_ACCOUNT_PROJECT/bexs-sandbox/g" ../.kubernetes/sandbox/serviceAccount.yaml

kubectl apply -f ../.kubernetes/sandbox/

set +x
