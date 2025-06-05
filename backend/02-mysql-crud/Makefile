# 📦 Makefile (fully containerized version)

SERVICE_NAME := app
MYSQL_SERVICE := db

build:
	docker compose build

up:
	docker compose up --detach

logs:
	docker compose logs -f $(SERVICE_NAME)

down:
	docker compose down

restart:
	docker compose down
	docker compose up --detach --build

db-logs:
	docker logs mysql-lab

db-cli:
	docker exec -it mysql-lab mysql -uuser -ppass learninglab

ps:
	docker compose ps

curl:
	curl http://localhost:8080/users

clean:
	docker compose down -v --remove-orphans

rebuild:
	make clean
	make up