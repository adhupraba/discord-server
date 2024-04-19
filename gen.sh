set -a            
source .env
set +a

url=$DB_URL
split=($(echo $url | tr "/" "\n"))
nameArr=($(echo "${split[2]}" | tr "?" "\n"))
dbname="${nameArr[0]}"
folder="internal/discord/public/model"

jet -source=postgres -dsn=$url -schema=public -path=internal/ -ignore-tables="goose_db_version" -ignore-enums="channel_type,member_role"
mv internal/$dbname internal/discord
sqlc generate
rm -rf sample.sql.go
mv db.go models.go $folder
find $folder -type f ! -name 'models.go' -exec rm -f {} \;