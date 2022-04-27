#!/bin/bash

curl -v -X POST http://localhost:8080/accounts/ \
   -H 'Content-Type: application/json' \
   -u "user:my_password" \
   -d '{"name":"localhost","url":"http://localhost","user_name":"u","password":"s"}'
