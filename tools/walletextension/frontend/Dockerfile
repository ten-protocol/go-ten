# Use an official Node.js LTS version as a base image
FROM node:20-alpine AS base

ARG NEXT_PUBLIC_NETWORK_NAME
ARG NEXT_PUBLIC_TENSCAN_URL
ARG NEXT_PUBLIC_GATEWAY_URL
ARG NEXT_PUBLIC_API_HOST_ENVIRONMENT
ARG NEXT_PUBLIC_CHAIN_ID

ENV NEXT_PUBLIC_NETWORK_NAME=$NEXT_PUBLIC_NETWORK_NAME
ENV NEXT_PUBLIC_TENSCAN_URL=$NEXT_PUBLIC_TENSCAN_URL
ENV NEXT_PUBLIC_GATEWAY_URL=$NEXT_PUBLIC_GATEWAY_URL
ENV NEXT_PUBLIC_API_HOST_ENVIRONMENT=${NEXT_PUBLIC_API_HOST_ENVIRONMENT}
ENV NEXT_PUBLIC_CHAIN_ID=$NEXT_PUBLIC_CHAIN_ID

# Set the working directory
WORKDIR /usr/src/app

# Copy the necessary files to the working directory
COPY tools/walletextension/frontend/ .

# Install dependencies
RUN npm ci

# Build the Next.js app
RUN npm run build

# Reduce the size of the final image by using a lighter base image
FROM node:20-alpine AS runner

# Set the working directory
WORKDIR /usr/src/app

# Copy only the necessary files from the build stage
COPY --from=base /usr/src/app/.next ./.next
COPY --from=base /usr/src/app/public ./public
COPY --from=base /usr/src/app/package*.json ./

# Install production dependencies
RUN npm ci --prod

# Set the environment variables
ENV PORT=80

# Expose the port
EXPOSE 80

# Start the application
CMD ["npm", "start"]
