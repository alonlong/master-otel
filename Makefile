all: stored ctld apid

stored:
	go build -o bin/stored cmd/stored/main.go

apid:
	go build -o bin/apid cmd/apid/main.go

ctld:
	go build -o bin/ctld cmd/ctld/main.go

pb:
	protoc --go_out=. --go_opt=module=master-otel --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,module=master-otel proto/stored/v1/*.proto
	protoc --go_out=. --go_opt=module=master-otel --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,module=master-otel proto/ctld/v1/*.proto

init-db:
	docker-compose kill postgres
	docker-compose rm -f postgres
	sudo rm -rf data
	docker-compose up -d