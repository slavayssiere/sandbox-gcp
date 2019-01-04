#!/bin/bash

PROJECT="cool-wharf-207907"

mkdir -p ./tpl
rm -f ./tpl/*

while IFS=';' read -r NAME_INJECTOR SIZE FREQUENCY REPLICAS
do
    cp ./injector-deploy-template.yaml ./tpl/injector-deploy-$NAME_INJECTOR.yaml

    cd ./tpl
        sed -i.bak s/"injector-deployment-name"/"injector-deployment-$NAME_INJECTOR"/g injector-deploy-$NAME_INJECTOR.yaml
        sed -i.bak s/"name_topic"/"projects\/$PROJECT\/topics\/${NAME_INJECTOR}_topic"/g injector-deploy-$NAME_INJECTOR.yaml
        sed -i.bak s/"MESSAGE_SIZE_TPL"/"$SIZE"/g injector-deploy-$NAME_INJECTOR.yaml
        sed -i.bak s/"FREQUENCY_PER_SECOND_TPL"/"$FREQUENCY"/g injector-deploy-$NAME_INJECTOR.yaml
        sed -i.bak s/"REPLICAS"/"$REPLICAS"/g injector-deploy-$NAME_INJECTOR.yaml
        
        rm -f injector-deploy-$NAME_INJECTOR.yaml.bak 
    cd ..

done < ../iac/list_injector.csv

kubectl create -f ./tpl/