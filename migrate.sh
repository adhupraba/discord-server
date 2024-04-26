set -a
source .env
set +a

goose -dir="internal/migrations" postgres $DB_URL up
