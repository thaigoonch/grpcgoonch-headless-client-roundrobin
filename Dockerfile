FROM golang:1.17 AS builder
WORKDIR /app
COPY . /app

ENV GOOS=linux

RUN chmod +x ./generate.sh && \
    ./generate.sh && \
    CGO_ENABLED=0 GOOS=linux \
    go install ./...
CMD ["/bin/go/grpcgoonchclient"]