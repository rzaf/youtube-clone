
all: generatePbs swagger build

generatePbs:
	@echo "generating porotocol buffer files ...."

	@mkdir -p ./database/pbs/helper/
	@cd ./database/pbs && protoc helper.proto --go_out=./helper --go_opt=paths=source_relative --go-grpc_out=./helper --go-grpc_opt=paths=source_relative  
	@mkdir -p ./database/pbs/comment/
	@cd ./database/pbs && protoc comment.proto --go_out=./comment --go_opt=paths=source_relative --go-grpc_out=./comment --go-grpc_opt=paths=source_relative  
	@mkdir -p ./database/pbs/media/
	@cd ./database/pbs && protoc media.proto --go_out=./media --go_opt=paths=source_relative --go-grpc_out=./media --go-grpc_opt=paths=source_relative  
	@mkdir -p ./database/pbs/playlist/
	@cd ./database/pbs && protoc playlist.proto --go_out=./playlist --go_opt=paths=source_relative --go-grpc_out=./playlist --go-grpc_opt=paths=source_relative  
	@mkdir -p ./database/pbs/user-pb/
	@cd ./database/pbs && protoc user.proto --go_out=./user-pb --go_opt=paths=source_relative --go-grpc_out=./user-pb --go-grpc_opt=paths=source_relative 

	@mkdir -p ./file/pbs/file/
	@cd ./file/pbs && protoc file.proto --go_out=./file --go_opt=paths=source_relative --go-grpc_out=./file --go-grpc_opt=paths=source_relative 

	@mkdir -p ./notification/pbs/notificationPb/
	@cd ./notification/pbs && protoc notification.proto --go_out=./notificationPb --go_opt=paths=source_relative --go-grpc_out=./notificationPb --go-grpc_opt=paths=source_relative 

	@mkdir -p ./email/pbs/emailPb/
	@cd ./email/pbs && protoc email.proto --go_out=./emailPb --go_opt=paths=source_relative --go-grpc_out=./emailPb --go-grpc_opt=paths=source_relative 

swagger:
	@echo "creating swagger docs of gateway service ...."

	@cd gateway && swag init -g ./handlers/docs.go
	@echo "creating swagger docs of auth service ...."
	@cd auth && swag init -g ./handlers/docs.go
	@echo "creating swagger docs of file service ...."
	@cd file && swag init -g ./handlers/docs.go

swagger-fmt:
	@echo "formatting swagger docs"
	@cd gateway/handlers && swag fmt
	@cd file/handlers && swag fmt

build-go:
	@echo "building go files:"
	
	@echo "building database service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o database/bin/migrateUp database/cmd/migrateUp/main.go
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o database/bin/migrateDown database/cmd/migrateDown/main.go
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o database/bin/databaseService database/cmd/databaseService/main.go

	@echo "building file service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o file/bin/fileService file/cmd/fileService/main.go

	@echo "building gateway service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gateway/bin/gatewayService gateway/cmd/gatewayService/main.go

	@echo "building email service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o email/bin/emailService email/cmd/emailService/main.go
	@cp -r email/email/templates email/bin

	@echo "building notification service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o notification/bin/notificationService notification/cmd/notificationService/main.go

	@echo "building authentication service ..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o auth/bin/authService auth/cmd/authService/main.go

build-docker:
	@echo "building docker containers ..."
	@docker compose build

build: build-go build-docker

run:
	@echo "running docker compose ..."
	docker compose up

test:
	@echo "running tests ..."
	@cd ./gateway/routes && go test -v .

remove:
	@echo "Stopping and removing containers ..."
	docker compose down
	
swarm: build-go
	@echo "building docker containers ..."
	@docker compose -f docker-compose.swarm.yml build

	@if [ "$$(docker info --format '{{.Swarm.LocalNodeState}}')" != "active" ]; then \
		echo "Initializing Docker Swarm..."; \
		docker swarm init || { echo "Failed to initialize Docker Swarm"; exit 1; }; \
	else \
		echo "Swarm is already active"; \
	fi

	@if [ -z "$$(docker secret ls | grep jwt_signing_key)" ]; then \
		echo "your-secret-key" | docker secret create jwt_signing_key -; \
	else \
		echo "Secret jwt_signing_key already exists"; \
	fi
	@docker stack deploy -c docker-compose.swarm.yml youtube-clone

watch:
	@echo "Watching for changes in Go files across services..." && \
	find ./gateway -name "*.go" | entr -r sh -c ' \
		echo "Changes in ./gateway detected, building..."; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gateway/bin/gatewayService gateway/cmd/gatewayService/main.go && echo "gateway/bin/gatewayService built successfully."' & \
	find ./file -name "*.go" | entr -r sh -c ' \
		echo "Changes in ./file detected, building..."; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o file/bin/fileService file/cmd/fileService/main.go && echo "file/bin/fileService built successfully."' & \
	find ./notification -name "*.go" | entr -r sh -c ' \
		echo "Changes in ./notification detected, building..."; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o notification/bin/notificationService notification/cmd/notificationService/main.go && echo "notification/bin/notificationService built successfully."' & \
	find ./database -name "*.go" | entr -r sh -c ' \
		echo "Changes in ./database detected, building..."; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o database/bin/databaseService database/cmd/databaseService/main.go && echo "database/bin/databaseService built successfully."' & \
	find ./auth -name "*.go" | entr -r sh -c ' \
		echo "Changes in ./auth detected, building..."; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o auth/bin/authService auth/cmd/authService/main.go && echo "auth/bin/authService built successfully."' & \
	wait

clean:
	@rm -f database/bin/*
	@rm -f notification/bin/*
	@rm -f file/bin/*
	@rm -f gateway/bin/*
	@docker compose down
