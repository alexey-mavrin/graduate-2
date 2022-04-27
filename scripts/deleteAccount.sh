#!/bin/bash

curl -v -X DELETE http://localhost:8080/accounts/1 \
   -H 'Content-Type: application/json' \
   -u "user:my_password"
