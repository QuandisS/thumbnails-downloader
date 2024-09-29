package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	downloader "thumbnails-downloader/pkg/downloader_v1"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	pflag.String("address", "localhost:9091", "Downloader server port")
	pflag.Bool("async", false, "Use async mode")
	pflag.Int("timeout", 10, "Timeout for requests in seconds")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	grpcClient, err := grpc.NewClient(viper.GetString("address"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Failed to create grpc client: ", err)
	}

	defer grpcClient.Close()

	client := downloader.NewDownloaderClient(grpcClient)

	wg := &sync.WaitGroup{}
	for i, url := range pflag.Args() {
		if viper.GetBool("async") {
			wg.Add(1)
			go func(url string, wg *sync.WaitGroup) {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("timeout"))*time.Second)
				defer cancel()
				r, err := client.Download(ctx, &downloader.DownloadRequest{
					Url: url})

				if err != nil {
					log.Fatal("Failed to download thumbnail: ", err)
				}

				outFile, err := os.Create(fmt.Sprintf("thumbnail_%d.jpg", i))
				if err != nil {
					log.Fatal(err)
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, bytes.NewReader(r.ImageData))

				if err != nil {
					log.Fatal(err)
				}

			}(url, wg)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("timeout"))*time.Second)
			defer cancel()
			r, err := client.Download(ctx, &downloader.DownloadRequest{
				Url: url})

			if err != nil {
				log.Fatal("Failed to download thumbnail: ", err)
			}

			outFile, err := os.Create(fmt.Sprintf("thumbnail_%d.jpg", i))
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, bytes.NewReader(r.ImageData))

			if err != nil {
				log.Fatal(err)
			}
		}
	}
	wg.Wait()
}
