#!/bin/bash
openssl genrsa -out proxy.key 2048
openssl req -new -key proxy.key -x509 -days 3650 -out proxy.crt -subj /C=CN/ST=BJ/O="Localhost Ltd"/CN=proxy