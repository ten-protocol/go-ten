#!/bin/bash

# Set the path to the docker-compose file
COMPOSE_FILE="./docker-compose.yaml"

# Clean up any existing containers and volumes
echo "== Removing old data..."
docker compose -f "${COMPOSE_FILE}" down -v

# Start the services in the correct order
echo "== Initializing Geth..."
docker compose -f "${COMPOSE_FILE}" up init_geth --wait

echo "== Creating beacon chain genesis..."
docker compose -f "${COMPOSE_FILE}" up create_beacon_chain_genesis --wait

echo "== Clearing beacon database..."
docker compose -f "${COMPOSE_FILE}" up clear_beacon_db --wait

echo "== Starting Geth..."
docker compose -f "${COMPOSE_FILE}" up -d geth
echo "== Waiting for Geth to be ready..."
sleep 2

echo "== Starting beacon chain..."
docker compose -f "${COMPOSE_FILE}" up -d prysm_beacon_chain
echo "== Waiting for beacon chain to be ready..."
sleep 2

echo "== Starting validator..."
docker compose -f "${COMPOSE_FILE}" up -d prysm_validator

echo "== All services started!"