#!/bin/bash

set -eo pipefail

go build -o fetch-locales-func
zip deployment.zip fetch-locales-func

aws lambda create-function \
--region eu-central-1 \
--function-name TranslationProxy \
--zip-file fileb://./deployment.zip \
--runtime go1.x \
--tracing-config Mode=Active \
--role arn:aws:iam::${AWS_ACCOUNT_ID}:role/PhraseProxyFull \
--handler main

rm fetch-locales-func deployment.zip