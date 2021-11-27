# Builder

FROM golang:1.17-alpine AS builder
RUN apk add --update make git curl

ARG GITHUB_PATH=github.com/ozonmp/ssn-service-api

COPY Makefile /home/${GITHUB_PATH}/Makefile
COPY go.mod /home/${GITHUB_PATH}/go.mod
COPY go.sum /home/${GITHUB_PATH}/go.sum
COPY pkg /home/${GITHUB_PATH}/pkg

WORKDIR /home/${GITHUB_PATH}

COPY . /home/${GITHUB_PATH}
RUN make build-go-bot

# Bot

FROM alpine:latest as server
ARG GITHUB_PATH=github.com/ozonmp/ssn-service-api
LABEL org.opencontainers.image.source=https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/bot .
COPY --from=builder /home/${GITHUB_PATH}/config.bot.yml .

RUN chown root:root bot

CMD ["./bot"]
