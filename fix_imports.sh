#!/bin/bash

# Find all Go files and update import paths
find . -name "*.go" -type f -exec sed -i '' 's|github.com/binomena/|github.com/igo-used/binomena/|g' {} \;

echo "Import paths updated successfully!"