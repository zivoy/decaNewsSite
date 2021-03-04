# todo make production and setip ci with github actions or something
FROM golang:1.16-alpine
MAINTAINER zivoy
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

ARG DB_CREDS
ARG D_KEY
ARG D_SECRET
ARG DEBUG_MODE=false
ARG HOST_PATH
# http://localhost:5000 is an example for testing
ARG RANDOM_SECRET

ENV DISCORD_KEY=$D_KEY
ENV DISCORD_SECRET=$D_SECRET
ENV DEBUG=$DEBUG_MODE
ENV REDIRECT=$HOST_PATH
ENV STORE_SECRET=$RANDOM_SECRET
ENV GOOGLE_APPLICATION_CREDENTIALS=$DB_CREDS

EXPOSE 5000

CMD ["./app"]