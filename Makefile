.PHONY: *

run:
	-docker compose -f ./deploy/local/run/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/docker-compose.yml up --build --attach=server

deploy-wa-scheduler-ec2:
	-docker-compose -f ./deploy/aws/ec2/docker-compose.yml down --remove-orphans
	docker-compose -f ./deploy/aws/ec2/docker-compose.yml up --build -d