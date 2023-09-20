docker-build:
	docker compose build
docker-up:
	docker compose up 
docker-bash:
	docker compose run --service-ports web bash 