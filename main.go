package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/devopshaven/static-site-service/miniohandler"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

var (
	listenAddr string
)

func init() {
	flag.StringVar(&listenAddr, "listen", "0.0.0.0:8080", "service listen address")
}

func main() {
	flag.Parse()

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_ACCESS_KEY_SECRET")

	useSSL := false

	siteName := os.Getenv("SITE_NAME")

	if strings.ToLower(os.Getenv("MINIO_USE_SSL")) == "true" {
		useSSL = true
	}

	// Initialize MinIO client object.
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Info().Err(err).Msg("cannot connect to minio instance")
	}

	miniohandler.PrometheusHandler()

	log.Info().Msgf("http server is listening on address: http://%s", listenAddr)
	if err := http.ListenAndServe(listenAddr, miniohandler.MinioHandler(mc, siteName, true)); err != nil {
		log.Warn().Err(fmt.Errorf("http server exited: %w", err)).Send()
	}
}
