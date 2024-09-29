package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	downloader "thumbnails-downloader/pkg/downloader_v1"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [--out <dir>] [--address <server-address>] [--async] [--timeout <timeout>] url...\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr, pflag.CommandLine.FlagUsages())
}

func main() {
	pflag.Usage = Usage
	pflag.String("address", "localhost:9091", "Downloader server port")
	pflag.Bool("async", false, "Use async mode")
	pflag.String("out", ".", "Output directory")
	pflag.Int("timeout", 10, "Timeout for requests in seconds")
	pflag.Parse()

	if pflag.CommandLine.NArg() < 1 {
		Usage()
		os.Exit(1)
	}

	viper.BindPFlags(pflag.CommandLine)

	outdir := viper.GetString("out")
	fstat, err := os.Stat(outdir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Output directory does not exist")
		}
		log.Fatal("Failed to access output directory: ", err)
	}
	if !fstat.IsDir() {
		log.Fatal("Output path is not directory")
	}

	grpcClient, err := grpc.NewClient(viper.GetString("address"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Failed to create grpc client: ", err)
	}

	defer grpcClient.Close()

	client := downloader.NewDownloaderClient(grpcClient)

	wg := &sync.WaitGroup{}
	for i, url := range pflag.Args() {
		filepath := path.Join(viper.GetString("out"), fmt.Sprintf("thumbnail_%d.jpg", i))
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

				outFile, err := os.Create(filepath)
				if err != nil {
					log.Fatal(err)
				}

				_, err = io.Copy(outFile, bytes.NewReader(r.ImageData))

				if err != nil {
					log.Fatal(err)
				}

				outFile.Close()
			}(url, wg)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("timeout"))*time.Second)
			defer cancel()
			r, err := client.Download(ctx, &downloader.DownloadRequest{
				Url: url})

			if err != nil {
				log.Fatal("Failed to download thumbnail: ", err)
			}

			outFile, err := os.Create(filepath)
			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(outFile, bytes.NewReader(r.ImageData))

			if err != nil {
				log.Fatal(err)
			}
			outFile.Close()
		}
	}
	wg.Wait()
}
