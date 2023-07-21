# Build stage for downloading dependencies based on the core defined system
FROM golang:1.20-buster as get-dependencies

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM get-dependencies as build

COPY . /home/obscuro/go-obscuro

# build the eth2network exec
WORKDIR /home/obscuro/go-obscuro/tools/obscuroscan_v2/backend/cmd
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o backend


FROM node:18-buster as runner
RUN mkdir -p /home/obscuro/go-obscuro/tools/obscuroscan_v2/backend/cmd
COPY --from=build /home/obscuro/go-obscuro/tools/obscuroscan_v2/ /home/obscuro/go-obscuro/tools/obscuroscan_v2/

WORKDIR /home/obscuro/go-obscuro/tools/obscuroscan_v2/frontend
RUN npm install && npm run build-prod && npm install http-server -g

WORKDIR /home/obscuro/go-obscuro/tools/obscuroscan_v2
EXPOSE 80