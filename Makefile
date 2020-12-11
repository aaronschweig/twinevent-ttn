.PHONY: dev
dev:
	go run main.go -config ./.default-config.yaml

.PHONY: docker-run
docker-run:
	docker run --network="host" -v $(PWD)/.default-config.yaml:/config/config.yaml:ro -d aaronschweig/twinevent-ttn -config config/config.yaml