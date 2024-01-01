SRC_APP=cmd/app/main.go
SRC_CONFIG=cmd/config/main.go
SRC_SENDER=cmd/sender/main.go

export CONFIG_PATH=./config/local.yaml

all: init_config
	go run $(SRC_APP)

init_config:
	go run $(SRC_CONFIG)

test: clean
	sudo docker-compose up

t:
	sudo docker exec -it server_container go run ./cmd/sender/main.go

tc: init_config
	sudo docker compose config

send:
	go run $(SRC_SENDER)

clean:
	sudo docker-compose down
	sudo docker container prune
	sudo rm -rf storage/pgdata
	rm -rf logs/logs.txt