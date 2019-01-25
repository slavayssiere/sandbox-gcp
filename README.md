# sandbox-gcp

## Objectif

Voir la présentation disponible [ici](perdu.com)

## Pré-requis

### A installer

- gcloud
- terraform

### Configuration

Créer le fichier de configuration dans le répertoire racine du projet.

```language-bash
#/bin/bash

export GCP_PROJECT="***"
export CONSUMER_KEY="***"
export CONSUMER_SECRET="***"
export ACCESS_TOKEN="***"
export ACCESS_SECRET="***"
```

## Création de la plateform

### Création de l'infrastructure

Aller dans le répertoire "iac". Puis lancer le script "create-all.sh".

### Déploiement

Aller dans le répertoire "deployment". Puis lancer le script "deploy.sh".

## Tests

### Flux temps réel

Aller sur le [site: "http://public.gcp-wescale.slavayssiere.fr"](http://public.gcp-wescale.slavayssiere.fr).

### API

#### Génération d'aggregas

```language-bash
curl -X POST http://public.gcp-wescale.slavayssiere.fr/aggregator/stats | jq .
```

```language-bash
curl -X GET http://public.gcp-wescale.slavayssiere.fr/aggregator/stats/1 | jq .
```

```language-bash
curl -X GET http://public.gcp-wescale.slavayssiere.fr/aggregator/top10 | jq .
```

#### To test datavisualisation

Aller sur Google Drive

kubectl create ns test
kubectl -n test run -i --tty busybox --image=busybox --restart=Never -- sh

## Observability

sudo ssh -i /Users/slavayssiere/.ssh/id_rsa admin@bastion.gcp-wescale.slavayssiere.fr -L 80:admin.gcp.wescale:80

http://servicegraph.localhost/force/forcegraph.html
http://jaeger.localhost
http://prometheus.localhost
http://grafana.localhost
