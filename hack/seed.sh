#! /bin/bash

DB=${1:-production}
USER=${2:-postgres}
HOST=${3:-localhost}
PORT=${4:-5432}

cat <<EOF | psql -U ${USER} -h ${HOST} -p ${PORT} -d ${DB}
INSERT INTO quotes (author, text)
VALUES ('Oscar Wilde', 'Be yourself; everyone else is already taken.');

INSERT INTO quotes (author, text) 
VALUES ('Friedrich Nietzsche', 'Without music, life would be a mistake.');

INSERT INTO quotes (author, text) 
VALUES ('Mahatma Gandhi', 'Be the change that you wish to see in the world.');

INSERT INTO quotes (author, text) 
VALUES ('Mark Twain', 'If you tell the truth, you don''t have to remember anything.');
EOF
