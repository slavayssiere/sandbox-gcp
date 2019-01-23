#!/bin/bash

# cbt -instance "test-instance" read "test-table"
bq --project_id "slavayssiere-sandbox" --location=EU mk test_bq
bq --project_id "slavayssiere-sandbox" --location=EU mk --external_table_definition=./simple_definition.json test_bq.ms

# gcloud ml language analyze-entities --content="Michelangelo Caravaggio, Italian painter, is known for 'The Calling of Saint Matthew'."