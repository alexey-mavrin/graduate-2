#!/bin/bash

curl -v -X GET http://localhost:8088/records/account/1 \
   -H 'Content-Type: application/json' \
   -u "user:my_password"
