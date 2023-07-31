#!/usr/bin/env bash
driver=$1
name=$2

eval test -n \"$name\" \
    || { echo "⛔ name is not set"; exit 1; }; \

echo "🦆 $1 check available goose..."
if ! command -v goose &> /dev/null
then
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi
echo "👍 goose [OK]"

echo "💾 creating DB migrations..."
goose -dir ./infrastructure/database/migration/$driver create $name sql
echo "🎉 DB migration finished!"