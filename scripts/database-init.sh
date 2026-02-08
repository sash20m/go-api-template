#!/bin/bash

set -euo pipefail

function create_user_and_database() {
	local database=$1
	echo "  Creating user and database '$database'"
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
	    CREATE USER $database;
	    CREATE DATABASE $database;
	    GRANT ALL PRIVILEGES ON DATABASE $database TO $database;
EOSQL
}

if [[ -z "${POSTGRES_DATABASE:-}" ]]; then
	echo "POSTGRES_DATABASE is not set; nothing to do"
	exit 0
fi

create_user_and_database "$POSTGRES_DATABASE"
