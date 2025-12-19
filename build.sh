#!/bin/bash
echo "Building The Invader..."
go build -o game ./cmd/invaders
if [ $? -eq 0 ]; then
    echo "Build successful! Run with ./game"
else
    echo "Build failed."
    exit 1
fi
