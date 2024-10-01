# Use an official Node.js 22 as a parent image
FROM node:22-alpine

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro/tools
COPY ./tools/bridge-frontend /home/obscuro/go-obscuro/tools/bridge-frontend

WORKDIR /home/obscuro/go-obscuro/tools/bridge-frontend

# ARG for build-time variable
ARG API_HOST
ARG API_HOST_ENVIRONMENT
ARG GOOGLE_ANALYTICS_ID
ARG INFURA_KEY

# ENV for URL to be used in the app
ENV NEXT_PUBLIC_BRIDGE_API_HOST=${API_HOST}
ENV NEXT_PUBLIC_API_HOST_ENVIRONMENT=${API_HOST_ENVIRONMENT}
ENV NEXT_PUBLIC_BRIDGE_GOOGLE_ANALYTICS_ID=${GOOGLE_ANALYTICS_ID}
ENV NEXT_PUBLIC_BRIDGE_INFURA_KEY=${INFURA_KEY}

RUN pnpm install

EXPOSE 80