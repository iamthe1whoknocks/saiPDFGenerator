up:
	docker-compose -f ./microservices/docker-compose.yml up -d

down:
	docker-compose -f ./microservices/docker-compose.yml down --remove-orphans

build:
	make service
	make docker

service:
		cd ./src/saiPDFGenerator/Boilerplate && go mod tidy && go build -o ./../../../microservices/saiPDFGenerator/build/sai-pdfgenerator
		cp ./src/saiPDFGenerator/Boilerplate/config.yml ./microservices/saiPDFGenerator/build/config.yml

docker:
	docker-compose -f ./microservices/docker-compose.yml up -d --build

log:
	docker-compose -f ./microservices/docker-compose.yml logs -f

logp:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-pdfgenerator
