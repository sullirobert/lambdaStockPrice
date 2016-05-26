#!/bin/bash

set -e
set -x

GOOS=linux go build -o main
# rm lambda.zip
zip -r -j lambda.zip main index.js
rm main

#API Gateway template marketplace
#
# {
#    "ticker": "$input.params('ticker')"
#}
