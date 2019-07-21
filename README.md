# Simple Admission Controller

[![Build Status](https://travis-ci.org/ChrisTheShark/simple-admission-controller.svg?branch=master)](https://travis-ci.org/ChrisTheShark/simple-admission-controller.svg?branch=master)

Simple example for creating a custom validating [admission controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) for Kubernetes.

## Getting Started

Clone this repository, customize the certificate authority configuration located at [certs/ca_config.txt](https://github.com/ChrisTheShark/simple-admission-controller/blob/master/certs/ca_config.txt) to include your location and domain information.

Execute gen_certs.sh to create certificates and install a secret in the cluster containing the certificate data for injection in your webhook. Next, execute `kubectl create -f manifest.yaml` at the root of this project. Now that the controller is running test both the compliant and non-compliant pod definitions located at the root of this project.
