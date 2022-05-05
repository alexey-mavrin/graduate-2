#!/bin/bash

curl -v -X POST http://localhost:8088/records/note \
   -H 'Content-Type: application/json' \
   -u "user:my_password" \
   -d '{"name":"note1","text":"note1 text"}'
