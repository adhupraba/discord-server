# Command to generate jet models

```bash
jet -source=postgres -dsn="postgresql://postgres:postgres@localhost:5432/discord?sslmode=disable" -schema=public -path=./.gen -ignore-tables="goose_db_version"
```

leveraging both `sqlc` and `go-jet` here

`sqlc` to convert migrations into models and `go-jet` to perform typesafe db queries

# Running development server

This automatically restarts dev server on code changes

```bash
compiledaemon --command="./discord-server"
```

# Build for production

```bash
go build -tags netgo -ldflags '-s -w' -o discord-server
```

# Create a goose migration file

```bash
goose create users sql
```

# Migration

use the migrate.sh to migrate the schema to the database

```bash
sh migrate.sh
```

# Generate typesafe models

To generate typesafe go models utilise the gen.sh helper script which uses sqlc to generate the models and go-jet to pull the schema from db to generate helper functions for the schema

```bash
sh gen.sh
```
