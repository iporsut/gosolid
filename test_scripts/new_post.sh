#!/usr/bin/env bash

curl -XPOST -H "Content-Type:application/json" "localhost:8080/posts"  -d '{
  "title": "'"$1"'",
  "body": "'"$2"'"
}'
