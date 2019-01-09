#!/bin/bash

PROJECT="cool-wharf-207907"

echo "test"
mkdir -p ./tpl
rm -Rf ./tpl/*

echo "test"

while IFS=';' read -r NAME_INJECTOR SIZE FREQUENCY
do
    echo "For $NAME_INJECTOR"
    
    echo "Create ${NAME_INJECTOR}_topic and ${NAME_INJECTOR}_subcription"
    gcloud alpha pubsub topics create "${NAME_INJECTOR}_topic"
    gcloud alpha pubsub subscriptions create "${NAME_INJECTOR}_subcription" \
      --ack-deadline 60 \
      --topic "${NAME_INJECTOR}_topic" \
      --topic-project $PROJECT

    cp ./normalizer-deploy-template.yaml ./tpl/normalizer-deploy-template-$NAME_INJECTOR.yaml

    cd ./tpl
        sed -i.bak s/"normalized-deployment-name"/"normalized-deployment-$NAME_INJECTOR"/g normalizer-deploy-template-$NAME_INJECTOR.yaml
        sed -i.bak s/"normalized-hpa-name"/"normalized-hpa-$NAME_INJECTOR"/g normalizer-deploy-template-$NAME_INJECTOR.yaml
        sed -i.bak s/"SUB_NAME_TPL"/"projects\/$PROJECT\/subscriptions\/${NAME_INJECTOR}_subcription"/g normalizer-deploy-template-$NAME_INJECTOR.yaml
        sed -i.bak s/"SUB_SHORT_NAME_TPL"/"${NAME_INJECTOR}_subcription"/g normalizer-deploy-template-$NAME_INJECTOR.yaml
        sed -i.bak s/"TOPIC_NAME_TPL"/"projects\/$PROJECT\/topics\/${NAME_INJECTOR}_topic"/g normalizer-deploy-template-$NAME_INJECTOR.yaml
        
        rm -f normalizer-deploy-template-$NAME_INJECTOR.yaml.bak 
    cd ..

done < ../iac/list_injector.csv

kubectl create -f ./tpl/