# Dockerfile for goldapps production
FROM golang:alpine AS buildStage
MAINTAINER digIT <digit@chalmers.it>

# Install git
RUN apk update
RUN apk upgrade
RUN apk add --update git

# Copy sources
RUN mkdir -p $GOPATH/src/github.com/cthit/goldapps
COPY . $GOPATH/src/github.com/cthit/goldapps
WORKDIR $GOPATH/src/github.com/cthit/goldapps/cmd

# Grab dependencies
RUN go get -d -v ./...

# build binary
RUN go install -v
RUN mkdir /app && mv $GOPATH/bin/cmd /app/goldapps

##########################
#    PRODUCTION STAGE    #
##########################
FROM alpine
MAINTAINER digIT <digit@chalmers.it>

# Add standard certificates
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Set user
RUN addgroup -S app
RUN adduser -S -G app -s /bin/bash app
USER app:app

# Copy execution script
COPY ./sleep_and_run.sh /app/sleep_and_run.sh

# Copy binary
COPY --from=buildStage /app/goldapps /app/goldapps

ENV WAIT 15s

# Set good defaults
WORKDIR /app
ENTRYPOINT ./sleep_and_run.sh
CMD -dry