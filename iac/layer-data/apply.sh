#!/bin/bash

REGION="europe-west1"

terraform apply \
    --var "region=$REGION" \
    -auto-approve


cbt -instance "test-instance" createfamily "test-table" ms