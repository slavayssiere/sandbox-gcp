#!/bin/bash

# cbt -instance "test-instance" read "test-table"
bq --project_id "slavayssiere-sandbox" --location=EU mk test_bq
bq --project_id "slavayssiere-sandbox" --location=EU mk --external_table_definition=./bt_test_definition.json test_bq.ms
