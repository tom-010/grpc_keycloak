#!/bin/bash

set -e

echo "See guide: https://github.com/grpc-up-and-running/samples/tree/master/ch06/token-based-authentication/certs"


host="localhost"

openssl genrsa -out server.key 2048
openssl req -new -x509 -days 365 -key server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out server.crt

openssl req -newkey rsa:2048 -nodes -key server.key -keyout server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=localhost" -out server.csr 
openssl x509 -req -extfile <(printf "subjectAltName=DNS:localhost,DNS:localhost") -days 365 -in server.csr -CA server.crt -CAkey server.key -CAcreateserial -out server.crt
