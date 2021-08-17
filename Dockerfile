FROM golang:1.16-alpine

LABEL maintainer="zivoy"

RUN apk add --no-cache git

WORKDIR /app/decafans-server

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY src src
COPY src/templates templates
COPY src/resources resources

RUN go build -o ./app ./src

# load the config from a local file rather then getting it from the web
ARG LOCAL_FILE=false

ARG VERSION

ENV DEV_MODE=$DEV_MODE
ENV VERSION=$VERSION

LABEL version=$VERSION

EXPOSE 5000

CMD ["./app"]