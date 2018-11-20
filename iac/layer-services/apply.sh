#!/bin/bash

cd ../layer-bastion
bastion=$(terraform output bastion-ip)
cd -

scp -r . admin@$bastion:~
