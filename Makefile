# Makefile for BirthdayMate

APP_NAME=birthdaymate
IMAGE_NAME=vinaypdb/$(APP_NAME)

build:
	go build -o bin/$(APP_NAME) main.go

run:
	go run main.go

docker:
	docker build -t $(IMAGE_NAME):latest .

push:
	docker push $(IMAGE_NAME):latest

tidy:
	go mod tidy

clean:
	rm -rf bin/

