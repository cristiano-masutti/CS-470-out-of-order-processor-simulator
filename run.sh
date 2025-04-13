#!/bin/bash
set -e

INPUT_FILE="$1"
OUTPUT_FILE="$2"

echo "🚀 Running simulator..."
./simulator "$INPUT_FILE" "$OUTPUT_FILE"
