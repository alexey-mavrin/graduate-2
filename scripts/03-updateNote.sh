#!/bin/bash

curl -v -X PUT http://localhost:8088/records/note/1 \
   -H 'Content-Type: application/json' \
   -u "user:my_password" \
   -d '{"name":"note1","text":"note 1 new text","meta":"x meta"}'
