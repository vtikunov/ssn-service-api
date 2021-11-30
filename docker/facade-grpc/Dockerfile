# Builder

FROM golang:1.17-alpine AS builder
RUN apk add --update make git curl

ARG GITHUB_PATH=github.com/ozonmp/ssn-service-api

COPY Makefile /home/${GITHUB_PATH}/Makefile
COPY go.mod /home/${GITHUB_PATH}/go.mod
COPY go.sum /home/${GITHUB_PATH}/go.sum
COPY pkg /home/${GITHUB_PATH}/pkg

WORKDIR /home/${GITHUB_PATH}
RUN make deps-go

COPY . /home/${GITHUB_PATH}
RUN make build-go-facade-grpc

# Facade-gRPC

FROM alpine:latest as server
ARG GITHUB_PATH=github.com/ozonmp/ssn-service-api
LABEL org.opencontainers.image.source=https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/facade-grpc .
COPY --from=builder /home/${GITHUB_PATH}/config.facade.grpc.yml .

RUN chown root:root facade-grpc

CMD ["./facade-grpc"]
