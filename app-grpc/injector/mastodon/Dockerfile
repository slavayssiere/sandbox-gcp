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
ENV SERVER https://linuxrocks.online
ENV CLIENT_ID 3615ede47f6fe8f0c41ab5110e4d26bceb70aafabc37f1a181c52ef25558fba8
ENV CLIENT_SECRET 3081bd6214bdd89250f3cc335b685cfc32ebb9760458415fe9dafdc10924139d
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp

# Localy do:
# export TOPIC_NAME="projects/sandbox-perrier/topics/raw-mastodon"
# export CLIENT_ID="924ceb63456ea64fb4ebcfaa6c78e50656cf41673fd0f97bcdae9eaf3f1b02be"
# export CLIENT_SECRET="50b373884cecc9b7141ec8b888c96b4b8c4c45671c617b44f0a659dd4efcb026"
# export SERVER="https://linuxrocks.online"
# export SECRET_PATH="/Users/perriea/Downloads/sandbox-perrier-9d47bea05833.json"
# export PROM_PORT=8080
# export LOGIN=
# export PASSWORD=