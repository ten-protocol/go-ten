#!/bin/sh
set -e

# Check if the CONFIG_YAML_BASE64 environment variable is set and not empty
if [ -n "$CONFIG_YAML_BASE64" ]; then
    echo "Decoding Config YAML configuration..."
    echo "$CONFIG_YAML_BASE64" | base64 -d > "./CONFIG_YAML_BASE64.yaml"
fi

# Check if the OVERRIDE_YAML_BASE64 environment variable is set and not empty
if [ -n "$OVERRIDE_YAML_BASE64" ]; then
    echo "Decoding Override YAML configuration..."
    echo "$OVERRIDE_YAML_BASE64" | base64 -d > "./OVERRIDE_YAML_BASE64.yaml"
fi

"$@"