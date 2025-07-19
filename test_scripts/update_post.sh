#!/usr/bin/env bash

curl -XPATCH -H "Content-Type:application/json" "localhost:8080/posts/$1"  -d '{
  "title": "'"$2"'",
  "body": "'"$3"'"
}'
