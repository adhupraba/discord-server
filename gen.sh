set -a            
source .env
set +a

url=$DB_URL
split=($(echo $url | tr "/" "\n"))
nameArr=($(echo "${split[2]}" | tr "?" "\n"))
dbname="${nameArr[0]}"
folder="internal/$dbname/public/model"

jet -source=postgres -dsn=$url -schema=public -path=internal -ignore-tables="goose_db_version"
sqlc generate
mv db.go models.go $folder
find $folder -type f ! -name 'models.go' -exec rm -f {} \;