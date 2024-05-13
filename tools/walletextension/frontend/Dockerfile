# Use an official Node.js 22 as a parent image
FROM node:22-alpine

WORKDIR /usr/src/app

# ARG for build-time variable (GATEWAY_API_URL)
ARG GATEWAY_API_URL

# ENV for URL to be used in the app
ENV NEXT_PUBLIC_API_GATEWAY_URL=${GATEWAY_API_URL}
ENV PORT=80

# Copy package.json and package-lock.json (or yarn.lock) into the container
COPY package*.json ./

RUN npm install
COPY . .
RUN npm run build
EXPOSE 80
CMD ["npm", "start"]