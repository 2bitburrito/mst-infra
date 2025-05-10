#!/bin/bash
echo "Currently not working"
echo "just use goose postgres <DBURL> up"

if [ -f .env.goose ]; then
  source .env.goose
else
  echo "Error: .env file not found. sukka"
  exit 1
fi

echo "GOOSE_DRIVER:" "$GOOSE_DRIVER"
echo "DB STRING:" "$GOOSE_DBSTRING"

goose "$GOOSE_DRIVER" "$GOOSE_DBSTRING" up
