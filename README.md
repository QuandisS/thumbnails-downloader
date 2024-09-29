# A YouTube thumbnails downloader
It consists of 3 elements:
1. gRPC proxy-server (accepts grpc call, makes get request to YT thumbnails service, responds with image data)
2. gRPC client
3. Redis (as a cache)

## Usage
**1.** Start a server and Redis:
```
cd deployments
docker compose -f deployments/docker-compose.yml --project-directory ./ up --build
```
Or you can use: `make start-server`  
If you want to do a clean start (force rebuild and anon volumes renew):
```
cd deployments
docker compose -f deployments/docker-compose.yml --project-directory ./ up --force-recreate --renew-anon-volumes --build
```
Or you can use: `make clean-start-server`

**2.** Do a request with grpc client, run `go run ./cmd/client/.` with parameters as it is described in "usage":
```
Usage: client [--out <dir>] [--address <server-address>] [--async] [--timeout <timeout>] url...
      --address string   Downloader server port (default "localhost:9091")
      --async            Use async mode
      --out string       Output directory (default ".")
      --timeout int      Timeout for requests in seconds (default 10)
```