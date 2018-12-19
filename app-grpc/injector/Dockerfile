FROM golang:alpine AS build-env
RUN apk add --no-cache git
ADD . /src
RUN cd /src && go get -v -d && go build -o goapp

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
# export FREQUENCY_PER_SECOND=4
# export MESSAGE_SIZE=1024 
# export TOPIC_NAME="projects/cool-wharf-207907/topics/test-topic"
# export SECRET_PATH="/Users/slavayssiere/go/src/github.com/challengerdeep/wescale-gcp-poc/iac/secret-sa-gcp-pubsub.json"
# export PROM_PORT=8080