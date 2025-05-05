.PHONY: *

run:
	-docker compose -f ./deploy/local/run/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/docker-compose.yml up --build --attach=server
