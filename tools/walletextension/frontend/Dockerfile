# Use an official Node.js LTS version as a base image
FROM node:20-alpine AS base

ARG GATEWAY_API_URL
ENV NEXT_PUBLIC_API_GATEWAY_URL=${GATEWAY_API_URL}

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
RUN npm ci --production


# Set the environment variables
ENV PORT=80

# Expose the port
EXPOSE 80

# Start the application
CMD ["npm", "start"]
