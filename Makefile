all: stored apid

stored:
	go build -o bin/stored cmd/stored/main.go

apid:
	go build -o bin/apid cmd/apid/main.go

pb:
	protoc --go_out=. --go_opt=module=master-otel --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,module=master-otel proto/common/v1/*.proto
	protoc --go_out=. --go_opt=module=master-otel --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,module=master-otel proto/stored/v1/*.proto
