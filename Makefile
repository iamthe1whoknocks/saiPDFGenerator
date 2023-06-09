up:
	docker-compose -f ./microservices/docker-compose.yml up -d

down:
	docker-compose -f ./microservices/docker-compose.yml down --remove-orphans

build:
	make service
	make docker

service:
		cd ./src/pdf-generator && go mod tidy && go build -o ../../microservices/pdf-generator/build/pdf-generator
		cp ./src/pdf-generator/config.yml ./microservices/pdf-generator/build/config.yml

docker:
	docker-compose -f ./microservices/docker-compose.yml up -d --build

log:
	docker-compose -f ./microservices/docker-compose.yml logs -f

logp:
	docker-compose -f ./microservices/docker-compose.yml logs -f pdf-generator
