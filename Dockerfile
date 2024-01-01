FROM golang:1.21

WORKDIR /app

COPY . .

ENV CONFIG_PATH "./config/local.yaml"

RUN go mod download && go mod verify

RUN go run ./cmd/config/main.go

RUN go build -o serverapp ./cmd/app/main.go

CMD [ "./serverapp" ]