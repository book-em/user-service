#!usr/bin/bash

openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private_key.key -out public_key.pem
chmod 400 private_key.key