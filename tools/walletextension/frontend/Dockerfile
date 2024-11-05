# Use an official Node.js LTS version as a base image
FROM node:20-alpine AS base

# Install pnpm
RUN npm install -g pnpm

ARG NEXT_PUBLIC_NETWORK_NAME
ARG NEXT_PUBLIC_TENSCAN_URL
ARG NEXT_PUBLIC_GATEWAY_URL
ARG NEXT_PUBLIC_API_HOST_ENVIRONMENT

ENV NEXT_PUBLIC_NETWORK_NAME=$NEXT_PUBLIC_NETWORK_NAME
ENV NEXT_PUBLIC_TENSCAN_URL=$NEXT_PUBLIC_TENSCAN_URL
ENV NEXT_PUBLIC_GATEWAY_URL=$NEXT_PUBLIC_GATEWAY_URL
ENV NEXT_PUBLIC_API_HOST_ENVIRONMENT=${NEXT_PUBLIC_API_HOST_ENVIRONMENT}

# Set the working directory
WORKDIR /usr/src/app

# Copy the necessary files to the working directory
COPY tools/walletextension/frontend/ .

# Install dependencies
RUN pnpm install --frozen-lockfile

# Build the Next.js app
RUN pnpm run build

# Reduce the size of the final image by using a lighter base image
FROM node:20-alpine AS runner

# Install pnpm
RUN npm install -g pnpm

# Set the working directory
WORKDIR /usr/src/app

# Copy only the necessary files from the build stage
COPY --from=base /usr/src/app/.next ./.next
COPY --from=base /usr/src/app/public ./public
COPY --from=base /usr/src/app/package*.json ./
COPY --from=base /usr/src/app/pnpm-lock.yaml ./

# Install production dependencies
RUN pnpm install --prod --frozen-lockfile

# Set the environment variables
ENV PORT=80

# Expose the port
EXPOSE 80

# Start the application
CMD ["pnpm", "start"]
