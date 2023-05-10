#!/usr/bin/env bash

#
# This script installs and runs datadog before running the actual command
# It requires that DD_API_KEY and DD_HOSTNAME is set, will fail otherwise
#

# 1. Check if the DD_API_KEY env var is set, exit with error otherwise
if [ -z "$DD_API_KEY" ] || [ -z "$DD_HOSTNAME" ]
then
    echo "Either DD_API_KEY and/or DD_HOSTNAME is not set. Exiting."
    exit 1
fi

sed -i "s/api_key: none/api_key: ${DD_API_KEY}/" /etc/datadog-agent/datadog.yaml
sed -i "s/hostname: none/hostname: ${DD_HOSTNAME}/" /etc/datadog-agent/datadog.yaml

service datadog-agent start

# 3. Execute whatever is in the arguments as a command
"$@"

