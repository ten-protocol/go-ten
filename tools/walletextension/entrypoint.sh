#!/bin/sh

# Start NGINX in the background
nginx &

# Start wallet_extension_linux with parameters passed to the script
/home/obscuro/go-obscuro/tools/walletextension/bin/wallet_extension_linux "$@"

# Wait for any process to exit
wait -n

# Exit with the status of the process that exited first
exit $?
