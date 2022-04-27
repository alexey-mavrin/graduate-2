#!/bin/bash

curl -v -X PUT http://localhost:8080/accounts/1 \
   -H 'Content-Type: application/json' \
   -u "user:my_password" \
   -d '{"name":"localhost","url":"http://localhost","user_name":"u","password":"x"}'
