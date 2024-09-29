package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"thumbnails-downloader/internal/downloader"
	dlproxy "thumbnails-downloader/pkg/downloader_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = 9091

type server struct {
	dlproxy.UnimplementedDownloaderServer
}

func (s *server) Download(ctx context.Context, req *dlproxy.DownloadRequest) (*dlproxy.DownloadResponse, error) {
	log.Println("Download: ", req.Url)

	imgData, err := downloader.DowloadThumbnail(req.Url)
	if err != nil {
		log.Printf("Failed to download thumbnail: %v", err)
		return nil, fmt.Errorf("failed to download thumbnail: %w", err)
	}

	return &dlproxy.DownloadResponse{
		ImageData: imgData,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	dlproxy.RegisterDownloaderServer(s, &server{})

	log.Println("Server listening on address: ", lis.Addr().String())

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
