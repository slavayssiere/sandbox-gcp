FROM golang:alpine AS build-env
RUN apk add --no-cache git
ADD . /src
RUN cd /src && go get -v -d && go build -o goapp

# final stage
FROM alpine
RUN apk add --no-cache --update ca-certificates
ENV FREQUENCY_PER_SECOND 4
ENV MESSAGE_SIZE 1024 
ENV SUB_NAME "projects/cool-wharf-207907/subscriptions/normalized_subcription"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
VOLUME /secret
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp

# Localy do:
# export SUB_NAME="projects/sandbox-perrier/subscriptions/raw-normalizer"
# export SECRET_PATH="/Users/perriea/Downloads/sandbox-perrier-9d47bea05833.json"
# export TOPIC_NAME="projects/sandbox-perrier/topics/raw-twitter"
# export PROM_PORT=8080

# docker build . -t wescale/consumer:1.0
# docker run -d -v $HOME/Downloads/sandbox-perrier-9d47bea05833.json:/secret/secret-sa-gcp-pubsub.json --env-file .env wescale/consumer:1.0