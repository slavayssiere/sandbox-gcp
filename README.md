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

#### To test datavisualisation

Aller sur Google Drive
