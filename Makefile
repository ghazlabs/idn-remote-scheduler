.PHONY: *

run:
	-docker compose -f ./deploy/local/run/docker-compose.yml down --remove-orphans
	docker compose -f ./deploy/local/run/docker-compose.yml up --build --attach=server-scheduler

deploy-ec2-wa-scheduler:
	-docker-compose -f ./deploy/aws/ec2/docker-compose.yml -p wa-scheduler down --remove-orphans
	docker-compose -f ./deploy/aws/ec2/docker-compose.yml -p wa-scheduler up --build -d