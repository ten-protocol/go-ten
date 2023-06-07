FROM golang:1.20-alpine

RUN apk add build-base

RUN mkdir /faucet
WORKDIR /faucet
COPY . .
RUN cd cmd && go build -o faucet .

EXPOSE 80