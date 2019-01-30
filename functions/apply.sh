#!/bin/bash

cd laststat
gcloud functions deploy laststat --entry-point LastStat --runtime go111 --trigger-http
cd -