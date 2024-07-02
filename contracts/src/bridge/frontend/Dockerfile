# Use an official Node.js 22 as a parent image
FROM node:22-alpine

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro/contracts/src/bridge
COPY ./contracts/src/bridge/frontend /home/obscuro/go-obscuro/contracts/src/bridge/frontend

WORKDIR /home/obscuro/go-obscuro/contracts/src/bridge/frontend

# ARG for build-time variable
ARG API_HOST
ARG L1_BRIDGE
ARG L2_BRIDGE
ARG MESSAGE_BUS
ARG GOOGLE_ANALYTICS_ID

# ENV for URL to be used in the app
ENV NEXT_PUBLIC_BRIDGE_API_HOST=${API_HOST}
ENV NEXT_PUBLIC_BRIDGE_L1=${L1_BRIDGE}
ENV NEXT_PUBLIC_BRIDGE_L2=${L2_BRIDGE}
ENV NEXT_PUBLIC_BRIDGE_MESSAGE_BUS=${MESSAGE_BUS}
ENV NEXT_PUBLIC_BRIDGE_GOOGLE_ANALYTICS_ID=${GOOGLE_ANALYTICS_ID}

RUN npm install

EXPOSE 80