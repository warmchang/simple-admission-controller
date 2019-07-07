# Simple Admission Controller

Simple example for creating a custom validating [admission controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) for Kubernetes.

## Getting Started

Clone this repository, customize the certificate authority configuration located at [certs/ca_config.txt](https://github.com/ChrisTheShark/simple-admission-controller/blob/master/certs/ca_config.txt) to include your location and domain information. Execute gen_certs.sh to create certificates and install a secret in the cluster containing the certificate data for injection in your webhook.
