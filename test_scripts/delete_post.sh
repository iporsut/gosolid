#!/usr/bin/env bash

curl -v -XDELETE "localhost:8080/posts/$1"
