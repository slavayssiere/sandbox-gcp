FROM golang:1.10-alpine AS build-env
RUN apk add --no-cache --update git
ADD . /src
RUN cd /src && go get -v -d && go build -o goapp

# final stage
FROM alpine
RUN apk add --no-cache --update ca-certificates
VOLUME /secret
ENV PROM_PORT 8080
ENV SUB_NAME "projects/slavayssiere-sandbox/subscriptions/test-sub"
ENV TOPIC_NAME "projects/slavayssiere-sandbox/topics/test-topic-2"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp

# Locally
# export SUB_NAME="projects/cool-wharf-207907/subscriptions/test-sub"
# export TOPIC_NAME="projects/cool-wharf-207907/topics/test-topic-2"
# export SECRET_PATH="/Users/slavayssiere/go/src/github.com/challengerdeep/wescale-gcp-poc/iac/secret-sa-gcp-pubsub.json"
# export PROM_PORT=8081
