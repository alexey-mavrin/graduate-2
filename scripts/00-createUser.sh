#!/bin/bash

curl -v -X POST http://localhost:8088/users \
   -H 'Content-Type: application/json' \
   -d '{"name":"user", "password":"my_password"}'