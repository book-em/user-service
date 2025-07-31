#!/bin/sh
set -e

echo "Waiting for the database..."
until pg_isready -h db -U bookem_userdb_user -d bookem_userdb_test; do
  sleep 2
done

echo "Ready! Running tests..."
go test -v ./test/integration/...