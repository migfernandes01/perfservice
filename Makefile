docker-build:
	docker compose build
docker-up:
	docker compose up 
docker-down:
	docker-compose down --volumes 
docker-bash:
	docker compose run --service-ports web bash 
jet-gen:
	jet -dsn='postgresql://postgres:postgres@localhost:5432/rinha?sslmode=disable' -sslmode=disable -schema=public -path=./gen
psql:
	docker exec -it  postgres psql -U postgres rinha