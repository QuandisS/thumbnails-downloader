package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"thumbnails-downloader/internal/downloader"
	dlproxy "thumbnails-downloader/pkg/downloader_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	dlproxy.UnimplementedDownloaderServer
	downloader *downloader.Downloader
}

func (s *server) Download(ctx context.Context, req *dlproxy.DownloadRequest) (*dlproxy.DownloadResponse, error) {
	log.Println("Download: ", req.Url)

	imgData, err := s.downloader.DowloadThumbnail(req.Url)
	if err != nil {
		log.Printf("Failed to download thumbnail: %v", err)
		return nil, fmt.Errorf("failed to download thumbnail: %w", err)
	}

	return &dlproxy.DownloadResponse{
		ImageData: imgData,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	server := &server{
		downloader: downloader.NewDownloader(),
	}

	dlproxy.RegisterDownloaderServer(s, server)

	log.Println("Server listening on address: ", lis.Addr().String())

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
