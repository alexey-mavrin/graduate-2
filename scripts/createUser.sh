#!/bin/bash

curl -v -X POST http://localhost:8080/users/ \
   -H 'Content-Type: application/json' \
   -d '{"name":"user6", "password":"my_password"}'
