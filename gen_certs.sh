#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [ "$#" -ne 2 ]; then
    echo "Please provide a service name and namespace." 
    echo "USAGE: gen_certs {SERVICE} {NAMESPACE}"
    exit 1
fi

export SERVICE=$1
export NAMESPACE=$2

# INJECT VARIABLES INTO THE CERT CONFIGURATION
sed -i .bak -e "s/{SERVICE}/$SERVICE/g" -e "s/{NAMESPACE}/$NAMESPACE/g" certs/simple_config.txt

# CREATE THE PRIVATE KEY FOR OUR CUSTOM CA
openssl genrsa -out certs/ca.key 2048

# GENERATE A CA CERT WITH THE PRIVATE KEY
openssl req -new -x509 -key certs/ca.key -out certs/ca.crt -config certs/ca_config.txt

# CREATE THE PRIVATE KEY FOR OUR SERVER
openssl genrsa -out certs/simple-key.pem 2048

# CREATE A CSR FROM THE CONFIGURATION FILE AND OUR PRIVATE KEY
openssl req -new -key certs/simple-key.pem -subj "/CN=$SERVICE.$NAMESPACE.svc" -out certs/simple.csr -config certs/simple_config.txt

# CREATE THE CERT SIGNING THE CSR WITH THE CA CREATED BEFORE
openssl x509 -req -in certs/simple.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/simple-crt.pem

# INJECT VARIABLES INTO THE WEBHOOK CONFIGURATION
CA_BUNDLE=$(cat certs/ca.crt | base64 | tr -d '\n')
sed -i .bak -e "s/{SERVICE}/$SERVICE/g" -e "s/{NAMESPACE}/$NAMESPACE/g" -e "s/{CA_BUNDLE}/$CA_BUNDLE/g" manifest.yaml

# CREATE SECRET CONTAINING CERT DATA
kubectl create secret generic $SERVICE -n $NAMESPACE \
      --from-file=key.pem=certs/simple-key.pem \
      --from-file=cert.pem=certs/simple-crt.pem
