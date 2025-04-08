#!/bin/bash
# For macOS
find . -type f -name "*.go" -exec sed -i '' 's|github.com/binomena/|github.com/igo-used/binomena/|g' {} \;
