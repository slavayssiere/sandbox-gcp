FROM golang:1.10-alpine AS build-env
RUN apk add --no-cache --update git
ADD . /src
RUN cd /src && go get -v -d && go build -o goapp

# final stage
FROM alpine
RUN apk add --no-cache --update ca-certificates
VOLUME /secret
ENV SUB_NAME "projects/slavayssiere-sandbox/subscriptions/test-sub"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp

# Locally
# export SUB_NAME="projects/cool-wharf-207907/subscriptions/test-sub"
# export SECRET_PATH="/Users/slavayssiere/go/src/github.com/challengerdeep/wescale-gcp-poc/iac/secret-sa-gcp-pubsub.json"
