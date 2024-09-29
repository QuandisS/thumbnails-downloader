PHONY: generate
generate:
	mkdir -p pkg/downloader_v1
	protoc --go_out=pkg/downloader_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/downloader_v1 --go-grpc_opt=paths=source_relative \
	api/downloader_v1/downloader.proto

start-server:
	cd deployments/
	docker compose -f deployments/docker-compose.yml --project-directory ./ up --force-recreate --renew-anon-volumes --build