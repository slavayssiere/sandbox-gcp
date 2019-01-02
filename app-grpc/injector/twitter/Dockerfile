FROM golang:alpine AS build-env
RUN apk add --no-cache git
ADD . /src
RUN cd /src && go build -o goapp

# final stage
FROM alpine
RUN apk add --no-cache --update ca-certificates
VOLUME /secret
ENV PROM_PORT 8080
ENV FREQUENCY_PER_SECOND 4
ENV MESSAGE_SIZE 1024 
ENV TOPIC_NAME "projects/slavayssiere-sandbox/topics/test-topic"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp

# Localy do:
# export TOPIC_NAME="projects/sandbox-perrier/topics/raw-twitter"
# export SECRET_PATH="/Users/perriea/Downloads/sandbox-perrier-9d47bea05833.json"
# export PROM_PORT=8080