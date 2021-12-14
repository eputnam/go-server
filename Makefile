TAG ?= latest
APP ?= health-check-server

.PHONY: docker

server-build:
	docker build . -t $(APP):$(TAG) --progress=plain

server-run:
	docker run -it --rm -p 8080:8080 --name $(APP) $(APP):$(TAG)

db-run:
	docker-compose up --detach db

db-stop:
	-docker-compose rm --stop --force db
