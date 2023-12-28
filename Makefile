SRC_APP=cmd/app/main.go
SRC_CONFIG=cmd/config/main.go

export CONFIG_PATH=./config/local.yaml

all: init_config

init_config:
	go run $(SRC_CONFIG)
	sudo docker compose config

test: clean
	sudo docker-compose up

clean:
	sudo docker-compose down
	sudo docker container prune
	sudo rm -rf storage/pgdata