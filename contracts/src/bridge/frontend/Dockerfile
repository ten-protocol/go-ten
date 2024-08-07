# Use an official Node.js 22 as a parent image
FROM node:22-alpine

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro/contracts/src/bridge
COPY ./contracts/src/bridge/frontend /home/obscuro/go-obscuro/contracts/src/bridge/frontend

WORKDIR /home/obscuro/go-obscuro/contracts/src/bridge/frontend

# ARG for build-time variable
ARG BRIDGE_ENVIRONMENT
ARG GOOGLE_ANALYTICS_ID

# ENV for URL to be used in the app
ENV NEXT_PUBLIC_BRIDGE_API_HOST_ENVIRONMENT=${BRIDGE_ENVIRONMENT}
ENV NEXT_PUBLIC_BRIDGE_GOOGLE_ANALYTICS_ID=${GOOGLE_ANALYTICS_ID}

RUN npm install

EXPOSE 80