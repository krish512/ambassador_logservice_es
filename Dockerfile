FROM golang:1.15 as build

WORKDIR /app

ADD . .

RUN go mod download

RUN go build -ldflags "-s -w" -o logservice_es

FROM ubuntu:20.10

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=build /app/logservice_es /usr/local/bin/logservice_es

ENTRYPOINT ["/usr/local/bin/logservice_es"]
