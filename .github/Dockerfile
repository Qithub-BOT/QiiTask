# =============================================================================
#  Test Container for Vaious Go Versions (ver. 2021.12.04-15:35)
# =============================================================================

# Default version
ARG VARIANT="1.15-alpine"

# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM golang:${VARIANT}

ENV GO111MODULE=on

RUN apk add --no-cache \
    build-base \
    alpine-sdk \
    git \
    bash

COPY ./go.mod /workspaces/go.mod
COPY ./go.sum /workspaces/go.sum

WORKDIR /workspaces

RUN go mod download

ENTRYPOINT go mod download && go test -race ./...
