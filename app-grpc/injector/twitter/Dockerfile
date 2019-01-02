FROM golang:alpine AS build-env
RUN apk add --no-cache git
ADD . /src
RUN cd /src && go get -v -d && go build -o goapp

# final stage
FROM alpine
RUN apk add --no-cache --update ca-certificates
VOLUME /secret
ENV PROM_PORT 8080
ENV TOPIC_NAME "projects/sandbox-perrier/topics/raw-twitter"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp

# Localy do:
# export TOPIC_NAME="projects/sandbox-perrier/topics/raw-twitter"
# export SECRET_PATH="/Users/perriea/Downloads/sandbox-perrier-9d47bea05833.json"
# export CONSUMER_KEY=""
# export CONSUMER_SECRET=""
# export ACCESS_TOKEN=""
# export ACCESS_SECRET=""
# export HASHTAG=""
# export PROM_PORT=8080

# docker build . -t wescale/injector-twitter:1.0
# docker run -d -v $HOME/Downloads/sandbox-perrier-9d47bea05833.json:/secret/secret-sa-gcp-pubsub.json --env-file .env wescale/injector-twitter:1.0