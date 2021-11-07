TAG ?= latest

.PHONY: docker

server-build:
	docker build . -t go-server:$(TAG) --progress=plain

server-run:
	docker run -it --rm -p 443:443 --name go-server go-server:$(TAG)

db-run:
	docker-compose up --detach db

db-stop:
	-docker-compose rm --stop --force db
