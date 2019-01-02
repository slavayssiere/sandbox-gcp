FROM golang:alpine AS build-env
RUN apk add --no-cache git
ADD . /src
RUN cd /src && go build -o goapp

# final stage
FROM alpine
RUN apk add --no-cache --update ca-certificates
VOLUME /secret
ENV PROM_PORT 8080
ENV MESSAGE_SIZE 1024 
ENV TOPIC_NAME "projects/sandbox-perrier/topics/raw-mastodon"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp

# Localy do:
# export TOPIC_NAME="projects/sandbox-perrier/topics/raw-mastodon"
# export CLIENT_ID=""
# export CLIENT_SECRET=""
# export SERVER="https://linuxrocks.online"
# export HASHTAG=""
# export SECRET_PATH="/Users/perriea/Downloads/sandbox-perrier-9d47bea05833.json"
# export PROM_PORT=8080
# export LOGIN=""
# export PASSWORD=""