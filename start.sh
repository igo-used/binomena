#!/bin/bash

#
# Copyright [2025] [uJ1NO (Juxhino Kapllanaj)] [binomena.com] [adaneural.com]
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# Print environment for debugging
echo "Environment variables:"
echo "PORT: $PORT"
echo "NODE_ID: $NODE_ID"
echo "DATA_DIR: $DATA_DIR"

# Create data directory if it doesn't exist
if [ -n "$DATA_DIR" ]; then
  echo "Setting up data directory: $DATA_DIR"
  mkdir -p $DATA_DIR
  
  # Create subdirectories for different types of data
  mkdir -p $DATA_DIR/blockchain
  mkdir -p $DATA_DIR/wallets
  
  echo "Data directory created and ready"
fi

# Start the blockchain node with supported flags
./binomena --api-port $PORT --p2p-port 9000 --id "$NODE_ID"