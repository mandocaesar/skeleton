#!/usr/bin/env bash

driver=$1
user=$2
password=$3
host=$4
dbname=$5
sslmode=$6

for var in driver user password host dbname ; do \
  eval test -n \"\$$var\" \
      || { echo "⛔ $var is not set"; exit 1; }; \
done

if [[ -z $sslmode ]]
then
    sslmode=
fi

echo "🦆 check available goose..."
if ! command -v goose &> /dev/null
then
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi
echo "👍 goose [OK]"

echo "💾 running DB migrations..."

if [[ $driver == "postgres" ]]
then
    goose -dir ./infrastructure/database/migration/postgres postgres "user=$user password=$password host=$host dbname=$dbname sslmode=$sslmode" up
    echo "🎉 DB postgres migration finished!"
elif [[ $driver == "mysql" ]]
then
    goose -dir ./infrastructure/database/migration/mysql mysql "$user:$password@tcp($host:3306)/$dbname" up
    echo "🎉 DB mysql migration finished!"
fi
