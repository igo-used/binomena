#!/bin/bash

# Print environment for debugging
echo "Environment variables:"
echo "PORT: $PORT"
echo "NODE_ID: $NODE_ID"
echo "DATA_DIR: $DATA_DIR"

# Create data directory if it doesn't exist
mkdir -p $DATA_DIR

# Check if we need to create a data directory symlink
# This helps the application find the data in its expected location
if [ -n "$DATA_DIR" ]; then
  echo "Setting up data directory: $DATA_DIR"
  
  # Create symbolic links if needed
  # This depends on how your application looks for data
  # You might need to adjust these paths based on your application
  ln -sf $DATA_DIR ./data 2>/dev/null || true
  
  # Export environment variable for the application
  export BINOMENA_DATA_DIR=$DATA_DIR
fi

# Start the blockchain node with supported flags
./binomena --api-port $PORT --p2p-port 9000 --id "$NODE_ID"