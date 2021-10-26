TAG ?= latest

.PHONY: docker

docker-build:
	docker build . -t go-server:$(TAG) --progress=plain

docker-run:
	docker run -it --rm -p 443:443 --name go-server go-server:$(TAG)
