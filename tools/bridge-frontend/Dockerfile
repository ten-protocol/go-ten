# Use an official Node.js 22 as a parent image
FROM node:22-alpine

# Install pnpm
RUN npm install -g pnpm

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro/tools
COPY ./tools/bridge-frontend /home/obscuro/go-obscuro/tools/bridge-frontend

WORKDIR /home/obscuro/go-obscuro/tools/bridge-frontend

# ARG for build-time variable
ARG NEXT_PUBLIC_BRIDGE_API_HOST
ARG NEXT_PUBLIC_API_HOST_ENVIRONMENT
ARG NEXT_PUBLIC_BRIDGE_GOOGLE_ANALYTICS_ID

# ENV for URL to be used in the app
ENV NEXT_PUBLIC_BRIDGE_API_HOST=${NEXT_PUBLIC_BRIDGE_API_HOST}
ENV NEXT_PUBLIC_API_HOST_ENVIRONMENT=${NEXT_PUBLIC_API_HOST_ENVIRONMENT}
ENV NEXT_PUBLIC_BRIDGE_GOOGLE_ANALYTICS_ID=${NEXT_PUBLIC_BRIDGE_GOOGLE_ANALYTICS_ID}

RUN pnpm install

EXPOSE 80
