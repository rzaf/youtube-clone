
all: generatePbs swagger build

generatePbs:
	@echo "generating porotocol buffer files ...."
	@cd .. && protoc youtube-clone/database/pbs/*.proto --go_out=. --go-grpc_out=.
	@cd .. && protoc youtube-clone/file/pbs/*.proto --go_out=. --go-grpc_out=.
	@cd .. && protoc youtube-clone/notification/pbs/*.proto --go_out=. --go-grpc_out=.

swagger:
	@echo "creating swagger docs of gateway service ...."
	@cd gateway && swag init -g ./handlers/docs.go
	@echo "creating swagger docs of file service ...."
	@cd file && swag init -g ./handlers/docs.go

swagger-fmt:
	@echo "formatting swagger docs"
	@cd gateway/handlers && swag fmt
	@cd file/handlers && swag fmt

build:
	@echo "building go files:"
	
	@echo "building database service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o database/bin/migrateUp database/cmd/migrateUp/main.go
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o database/bin/migrateDown database/cmd/migrateDown/main.go
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o database/bin/databaseService database/cmd/databaseService/main.go

	@echo "building file service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o file/bin/fileService file/cmd/fileService/main.go

	@echo "building gateway service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gateway/bin/gatewayService gateway/cmd/gatewayService/main.go

	@echo "building notification service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o notification/bin/notificationService notification/cmd/notificationService/main.go


build-run: all
	@echo "building and running docker-compose ..."
	docker-compose up --build

run:
	@echo "running docker-compose ..."
	docker-compose up --build


remove:
	@echo "Stopping and removing containers ..."
	docker-compose down
	

clean:
	@rm -f database/bin/*
	@rm -f notification/bin/*
	@rm -f file/bin/*
	@rm -f gateway/bin/*
	@docker-compose down
