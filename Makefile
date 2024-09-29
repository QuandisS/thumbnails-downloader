PHONY: generate
generate:
	mkdir -p pkg/downloader_v1
	protoc --go_out=pkg/downloader_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/downloader_v1 --go-grpc_opt=paths=source_relative \
	api/downloader_v1/downloader.proto