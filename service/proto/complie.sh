protoc --proto_path=./proto --go_out=./ ./proto/*.proto
protoc --proto_path=./proto --go_out=plugins=grpc:. -I=./proto -I=. common.proto