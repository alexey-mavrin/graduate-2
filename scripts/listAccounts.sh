#!/bin/bash

curl -v -X GET http://localhost:8080/accounts \
   -H 'Content-Type: application/json' \
   -u "user:my_password"
