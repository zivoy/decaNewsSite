FROM golang:1.16 as BUILDER
LABEL stage=builder

WORKDIR /app/decafans-server

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY src src

RUN CGO_ENABLED=0 GOOS=linux go build -o ../main ./src

RUN mv ./src/templates ../templates && mv ./src/resources ../resources


FROM alpine:latest
LABEL maintainer="zivoy"

WORKDIR /root/

COPY --from=BUILDER /app/main app
COPY --from=BUILDER /app/templates templates
COPY --from=BUILDER /app/resources resources

# load the config from a local file rather then getting it from the web
ARG LOCAL_FILE=false

ARG VERSION

ENV DEV_MODE=$DEV_MODE
ENV VERSION=$VERSION

LABEL version=$VERSION

EXPOSE 5000

CMD ["./app"]