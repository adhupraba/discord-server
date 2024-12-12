The UI for this project is live at: https://discord.adhupraba.com

The source code for the UI is at: https://github.com/adhupraba/discord-client

---

# Command to generate jet models

```bash
jet -source=postgres -dsn="postgresql://postgres:postgres@localhost:5432/discord?sslmode=disable" -schema=public -path=./.gen -ignore-tables="goose_db_version"
```

leveraging both `sqlc` and `go-jet` here

`sqlc` to convert migrations into models and `go-jet` to perform typesafe db queries

# Running development server

[compiledaemon](https://github.com/githubnemo/CompileDaemon) automatically restarts dev server on code changes

```bash
compiledaemon --command="./discord-server"
```

or you can use [air](https://github.com/cosmtrek/air) to restart dev server

```bash
air
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

# Build for production

```bash
go build -tags netgo -ldflags '-s -w' -o discord-server
```

1. **`go build`**:
   This is the Go command to compile Go programs. By default, it compiles the Go code in the current directory (unless a specific package or file is mentioned).

2. **`-tags netgo`**:
   Build tags are a way to include/exclude certain files from the build process based on conditions.

   - `netgo`: This specific tag tells the Go compiler to use the Go-based `net` package instead of the system's native networking libraries. This can be helpful in situations where you want your application to be purely Go-based without relying on system C libraries for networking, potentially increasing portability and reducing issues with library dependencies.

3. **`-ldflags '-s -w'`**:
   The `-ldflags` (linker flags) option allows you to pass flags to the Go linker.

   - `-s`: Omit the symbol table and debug information.
   - `-w`: Omit the DWARF symbol table.

   Using these flags (`-s` and `-w`) reduces the size of the resulting binary by removing symbol tables and debug information. It's a common practice to use these flags when building release versions of applications where you want to minimize the binary size and aren't concerned with having debug information embedded in the binary.

4. **`-o app`**:
   The `-o` flag specifies the output file name for the compiled binary. In this case, the resulting binary will be named `app`.

In summary, this command compiles the Go code in the current directory using the Go-based `net` package, omits symbol and debug information to reduce binary size, and outputs the resulting binary with the name `app`.
