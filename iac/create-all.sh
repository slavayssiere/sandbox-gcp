#!/bin/bash

cd layer-base
./apply.sh
cd -

cd layer-kubernetes
./apply.sh
cd -