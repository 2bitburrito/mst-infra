#!/bin/bash

if [ -f .env.goose ]; then
  source .env.goose
else
  echo "Error: .env file not found."
  exit 1
fi

echo "GOOSE_DRIVER:" "$GOOSE_DRIVER"
echo "DB STRING:" "$GOOSE_DBSTRING"

cd migrations || exit

goose "$GOOSE_DRIVER" "$GOOSE_DBSTRING" down
