#! /bin/bash

DB=${1:-production}
USER=${2:-postgres}
HOST=${3:-localhost}
PORT=${4:-5432}

echo "DELETE FROM quotes" | psql -U ${USER} -h ${HOST} -p ${PORT} -d ${DB}
