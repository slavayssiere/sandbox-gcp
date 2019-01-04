FROM golang:1.10-alpine AS build-env
RUN apk add --no-cache --update git
ADD . /src
RUN cd /src && go get -v -d && go build -o goapp

# final stage
FROM alpine
RUN apk add --no-cache --update ca-certificates
VOLUME /secret
ENV PROM_PORT 8080
ENV SUB_NAME "projects/sandbox-perrier/topics/raw-twitter"
ENV TOPIC_NAME "projects/sandbox-perrier/topics/raw-normalizer"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp

# Locally
# export SUB_NAME="projects/sandbox-perrier/subscriptions/raw-twitter"
# export TOPIC_NAME="projects/sandbox-perrier/topics/raw-normalizer"
# export SECRET_PATH="/Users/perriea/Downloads/sandbox-perrier-9d47bea05833.json"
# export PROM_PORT=8080
