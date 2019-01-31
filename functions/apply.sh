#!/bin/bash

cd laststat
gcloud functions deploy laststat \
    --entry-point LastStat \
    --runtime go111 \
    --trigger-http \
    --service-account service-500721978414@gcf-admin-robot.iam.gserviceaccount.com
cd -

cd getstat
gcloud functions deploy getstat \
    --entry-point GetStat \
    --runtime go111 \
    --trigger-http \
    --service-account service-500721978414@gcf-admin-robot.iam.gserviceaccount.com
cd -

