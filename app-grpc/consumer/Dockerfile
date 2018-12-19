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