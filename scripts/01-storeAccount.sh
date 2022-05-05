#!/bin/bash

curl -v -X POST http://localhost:8088/records/account \
   -H 'Content-Type: application/json' \
   -u "user:my_password" \
   -d '{"name":"localhost","url":"http://localhost","user_name":"u","password":"s"}'
