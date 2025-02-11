#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER shdgfhjd WITH PASSWORD 'skdjfskd';
    CREATE DATABASE sdfsdfsdf;
    GRANT ALL PRIVILEGES ON DATABASE sdfsdfsdf TO shdgfhjd;
EOSQL

echo "user and database created"