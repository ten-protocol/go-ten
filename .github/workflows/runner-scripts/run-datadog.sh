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

# Install datadog agent
DD_API_KEY=none DD_HOSTNAME=none DD_INSTALL_ONLY=true DD_SITE="datadoghq.eu" bash -c "$(curl -L https://s3.amazonaws.com/dd-agent/scripts/install_script_agent7.sh)"
sed -i 's/# logs_enabled: false/logs_enabled: true/' /etc/datadog-agent/datadog.yaml
mkdir -p /etc/datadog-agent/conf.d/eth2network.d/
cp -p /home/obscuro/go-obscuro/.github/workflows/runner-scripts/conf.yml /etc/datadog-agent/conf.d/eth2network.d/conf.yml

sed -i "s/api_key: none/api_key: ${DD_API_KEY}/" /etc/datadog-agent/datadog.yaml
sed -i "s/hostname: none/hostname: ${DD_HOSTNAME}/" /etc/datadog-agent/datadog.yaml

service datadog-agent start

# 3. Execute whatever is in the arguments as a command
"$@"

