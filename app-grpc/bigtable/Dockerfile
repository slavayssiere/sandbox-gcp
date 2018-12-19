FROM golang:1.10-alpine AS build-env
RUN apk add --no-cache --update git
ADD . /src
RUN cd /src && go get -v -d && go build -o goapp

# final stage
FROM alpine
RUN apk add --no-cache --update ca-certificates
VOLUME /secret
ENV PROJECT_ID "cool-wharf-207907"
ENV TABLE_ID "my-table"
ENV INSTANCE_ID "test-1"
ENV SUB_NAME "projects/cool-wharf-207907/subscriptions/test-sub"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
ENV PROM_PORT 8080
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp


# Localy do:
# export PROJECT_ID="cool-wharf-207907"
# export TABLE_ID="my-table"
# export INSTANCE_ID="bt-instance-test"
# export SUB_NAME="projects/cool-wharf-207907/subscriptions/trade_subcription"
# export SECRET_PATH="/Users/slavayssiere/go/src/github.com/challengerdeep/wescale-gcp-poc/iac/secret-sa-gcp-pubsub.json"
# export PROM_PORT=8082