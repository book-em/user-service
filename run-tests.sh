#!/bin/bash

set -euo pipefail

echo "Running tests..."

cd src
go test -coverprofile=coverage.out -coverpkg=./... ./...
go tool cover -func=coverage.out

echo "Linting..."

go vet ./...

echo "All checks passed!"