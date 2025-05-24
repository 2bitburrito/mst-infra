#!/bin/bash

#goose postgres "postgres://mst_admin:Q70AqiE8KOfRIHxmqmN4@tf-20250502141116491500000001.cvq42ycqkt4f.us-west-1.rds.amazonaws.com:5432/mst_db" up

if [ -f .env.goose ]; then
  source .env.goose
else
  echo "Error: .env file not found."
  exit 1
fi

echo "GOOSE_DRIVER:" "$GOOSE_DRIVER"
echo "DB STRING:" "$GOOSE_DBSTRING"

cd migration || exit

goose "$GOOSE_DRIVER" "$GOOSE_DBSTRING" up
