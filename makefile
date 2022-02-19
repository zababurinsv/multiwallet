install:
		cd cmd/multiwallet && go install

protos:
		cd api/pb && protoc --go_out=paths=source_relative:. api.proto