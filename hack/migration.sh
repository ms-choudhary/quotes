#! /bin/bash

DB=${1:-production}
USER=${2:-postgres}
HOST=${3:-localhost}
PORT=${4:-5432}

echo "CREATE DATABASE ${DB};" | psql -U ${USER} -h ${HOST} -p ${PORT} || true

cat <<EOF | psql -U ${USER} -h ${HOST} -p ${PORT} -d ${DB} || true
CREATE TABLE quotes (
  id SERIAL PRIMARY KEY,
  text VARCHAR(256) UNIQUE,
  author VARCHAR(50)
);
EOF
