# evolve_auth
Go Backend for auth.

# Dev Setup

## Start CockroachDB

1. Install CockroachDB
2. Run the below command to start a single node cockroachdb instance.
```sh
cockroach start-single-node --insecure --listen-addr=localhost:26257 --http-addr=localhost:8080
```

## Start the server

1. Install Go
2. Set Environment Variables. Execute the following command in the root directory of the project.
```.env
export DATABASE_URL=<database_url>
export MAILER_EMAIL=<mailer_email>
export MAILER_PASSWORD=<mailer_password>
export FRONTEND_URL=<frontend_url>
export HTTP_PORT=5000
export GRPC_PORT=5001
export INIT_DB=true
```
3. Install the protobuf-grpc compiler.
```sh
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```
4. Run the following command to start the server.
```sh
go run main.go
```

# Editing .proto files

1. Install protoc compiler
2. Run the following command to generate the go files from the proto files.
```sh
protoc --go_out=./ --go_opt=paths=source_relative \
    --go-grpc_out=./ --go-grpc_opt=paths=source_relative \
    ./proto/authenticate.proto
```
