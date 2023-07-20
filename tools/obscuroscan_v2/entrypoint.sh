#!/bin/bash

# turn on bash's job control
set -m

# Start the primary process and put it in the background
./backend/cmd/backend --nodeHostAddress $NODEHOSTADDRESS --serverAddress 0.0.0.0:8080 &

# Start the helper process
http-server ./frontend/dist/ -a 0.0.0.0 -p 80 &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?