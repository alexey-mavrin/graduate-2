#!/bin/bash

curl -v -X GET http://localhost:8088/ping \
   -H 'Content-Type: application/json' \
   -u "user:my_password"
