#!/bin/bash

apps=("api" "cli" "worker")
values=("api" "drive" "worker")

for i in "${!apps[@]}"; do
  echo "Building ${values[$i]}..."
  go build -o  bin/${values[$i]} ./cmd/${apps[$i]}/main.go
done
