docker-build:
	docker compose build
docker-up:
	docker compose up 
docker-down:
	docker-compose down --volumes 
docker-bash:
	docker compose run --service-ports web bash 