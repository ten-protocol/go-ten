# Node.js LTS version
FROM node:18-buster AS base

# Install pnpm
RUN npm install -g pnpm

# build-time arguments
ARG NEXT_PUBLIC_NETWORK_NAME
ARG NEXT_PUBLIC_TENSCAN_URL
ARG NEXT_PUBLIC_GATEWAY_URL
ARG NEXT_PUBLIC_API_HOST_ENVIRONMENT

# environment variabless
ENV NEXT_PUBLIC_NETWORK_NAME=$NEXT_PUBLIC_NETWORK_NAME
ENV NEXT_PUBLIC_TENSCAN_URL=$NEXT_PUBLIC_TENSCAN_URL
ENV NEXT_PUBLIC_GATEWAY_URL=$NEXT_PUBLIC_GATEWAY_URL
ENV NEXT_PUBLIC_API_HOST_ENVIRONMENT=$NEXT_PUBLIC_API_HOST_ENVIRONMENT

# Set the working directory
WORKDIR /usr/src/app

# copy app source
COPY tools/walletextension/frontend/ .

# Install dependencies
RUN pnpm install

# Build the Next.js app
RUN pnpm run build

# using a lighter image for the final build
FROM node:18-buster AS runner

# working directory
WORKDIR /usr/src/app

# Copy the necessary files from the base image
COPY --from=base /usr/src/app/.next ./.next
COPY --from=base /usr/src/app/public ./public
COPY --from=base /usr/src/app/package*.json ./

# Install only production dependencies
RUN pnpm install --production

# Set environment variable and expose port
ENV PORT=80
EXPOSE 80

# start the app
CMD ["pnpm", "start"]
