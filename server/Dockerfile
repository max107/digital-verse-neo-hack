FROM golang:1.16-alpine

ENV BLOOM_PORT 8080
ENV APP_NAME DigitalVerse

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

ADD ./main.go  /go/src/${APP_NAME}/
ADD ./config.go  /go/src/${APP_NAME}/
ADD ./s3.go  /go/src/${APP_NAME}/
ADD ./neofs-cli  /go/src/${APP_NAME}/

RUN chmod +x /go/src/${APP_NAME}/neofs-cli

WORKDIR /go/src/${APP_NAME}/

RUN go mod init

RUN go get ./...

RUN go install 

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/${APP_NAME}

# Document that the service listens on port 8080.
EXPOSE ${BLOOM_PORT}