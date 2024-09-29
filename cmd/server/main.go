package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"thumbnails-downloader/pkg/downloader_v1/api/downloader_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = 9091

type server struct {
	downloader_v1.UnimplementedDownloaderServer
}

func (s *server) Download(ctx context.Context, req *downloader_v1.DownloadRequest) (*downloader_v1.DownloadResponse, error) {
	log.Println("Download: ", req.Url)

	return &downloader_v1.DownloadResponse{
		ImageData: []byte("image data"),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	downloader_v1.RegisterDownloaderServer(s, &server{})

	log.Println("Server listening on address: ", lis.Addr().String())

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
