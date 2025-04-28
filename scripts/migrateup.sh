#!/bin/bash
# THIS ISN"T ACTUALLY CONNECTED TO ANYTHING RN

if [ -f .env ]; then
    source .env
fi

cd db

goose posttgres $DATABASE_URL up