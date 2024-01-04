SRC_APP=cmd/app/main.go
SRC_CONFIG=cmd/config/main.go
SRC_SENDER=cmd/sender/main.go
SERVER_IMAGE_NAME=my_server_image

export CONFIG_PATH=./config/local.yaml

all: stop
	cd scripts && ./server_open.sh &
	sudo docker-compose up

depend:
	go mod download
	go mod verify

clean_start: clean
	cd scripts && ./server_open.sh &
	sudo docker-compose up

clean_start_with_send: clean
	cd scripts && ./server_open.sh &
	sudo docker-compose up
	sudo docker exec -it server_container go run ./cmd/sender/main.go

silent_start_with_send: stop
	cd scripts && ./server_open.sh &
	sudo docker-compose up -d
	sudo docker exec -it server_container go run ./cmd/sender/main.go

server: init_config
	go run $(SRC_APP)

send: init_config
	go run $(SRC_SENDER)

send_docker:
	sudo docker exec -it server_container go run ./cmd/sender/main.go

init_config:
	go run $(SRC_CONFIG)

clean_ports: init_config
	cd scripts && ./cleanports.sh

clean:
	sudo docker-compose down
	sudo rm -rf storage/pgdata
	rm -rf logs/logs.txt

stop:
	sudo docker-compose stop

cleanall: clean
	sudo docker rmi $(SERVER_IMAGE_NAME)