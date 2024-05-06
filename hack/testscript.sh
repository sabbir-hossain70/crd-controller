#!/bin/bash

OUTPUT_DIR="./manifests"
ABSOLUTE_PATH=$(realpath "$OUTPUT_DIR")

# Print the current working directory
echo "Current working directory is: $(pwd)"

# Print the resolved absolute path
echo "Resolved absolute path is: $ABSOLUTE_PATH"

# Check if the directory exists
if [[ -d "$ABSOLUTE_PATH" ]]; then
  echo "Directory exists."
else
  echo "Directory does not exist. Creating it now..."
  mkdir -p "$ABSOLUTE_PATH"
fi

# Now run your controller-gen command
controller-gen rbac:roleName=my-crd-controller crd \
  paths=github.com/sabbir-hossain70/crd/pkg/apis/crd.com/v1 \
  output:crd:dir="$ABSOLUTE_PATH" output:stdout
