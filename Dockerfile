FROM golang:1.16 as BUILDER
LABEL stage=builder

WORKDIR /app/decafans-server

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY src src

RUN CGO_ENABLED=0 GOOS=linux go build -o ../main ./src

RUN mv ./src/templates ../templates && mv ./src/resources ../resources


FROM node:alpine as COMPILER

WORKDIR /app
RUN npm install --global sass

COPY sass sass
COPY package.json package.json
COPY package-lock.json package-lock.json

RUN mkdir css

RUN npm install
RUN sass --no-source-map sass/main.sass css/main.css --style compressed
RUN sass --no-source-map sass/submitter.sass css/poster.css --style compressed

FROM alpine:latest
LABEL maintainer="zivoy"

WORKDIR /root/

COPY --from=BUILDER /app/main app
COPY --from=BUILDER /app/templates templates
COPY --from=BUILDER /app/resources resources

COPY --from=COMPILER /app/css resources/css

# load the config from a local file rather then getting it from the web
ARG LOCAL_FILE=false

ARG VERSION

ENV DEV_MODE=$DEV_MODE
ENV VERSION=$VERSION

LABEL version=$VERSION

EXPOSE 5000

CMD ["./app"]