# Build stage for downloading dependencies based on the core defined system
FROM golang:1.25-bookworm as get-dependencies

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM get-dependencies as build

COPY . /home/obscuro/go-obscuro

# build exec
WORKDIR /home/obscuro/go-obscuro/tools/tenscan/backend/cmd
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o backend

WORKDIR /home/obscuro/go-obscuro/tools/tenscan/backend
EXPOSE 80