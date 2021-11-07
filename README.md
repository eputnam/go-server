# Go webserver
I'm making a webserver in Go!

# Run the server
Run this command in the project root
```shell
go run ./server.go
```

# Docker

## The server
Build the Docker image (optionally set TAG)
```shell
TAG=my-tag make server-build
```

Run the Docker container (optionally set TAG)
```shell
TAG=my-tag make server-run
```

## PostgreSQL
You can run a local PostgreSQL server with these make targets:
```shell
make db-run
make db-stop
```