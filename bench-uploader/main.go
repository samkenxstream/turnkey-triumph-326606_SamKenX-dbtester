package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"

	"github.com/spf13/cobra"
)

// TODO: vendor gRPC and combine this into bench
// Currently, we get:
// panic: http: multiple registrations for /debug/requests

var Command = &cobra.Command{
	Use:   "bench-uploader",
	Short: "Uploads to cloud storage.",
	RunE:  CommandFunc,
}

var (
	csvResultPath                 string
	googleCloudProjectName        string
	googleCloudStorageJSONKeyPath string
	googleCloudStorageBucketName  string
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	Command.PersistentFlags().StringVar(&csvResultPath, "csv-result-path", "timeseries.csv", "path to store csv results.")
	Command.PersistentFlags().StringVar(&googleCloudProjectName, "google-cloud-project-name", "", "Google cloud project name.")
	Command.PersistentFlags().StringVar(&googleCloudStorageJSONKeyPath, "google-cloud-storage-json-key-path", "", "Path of JSON key file.")
	Command.PersistentFlags().StringVar(&googleCloudStorageBucketName, "google-cloud-storage-bucket-name", "", "Google cloud storage bucket name.")
}

func main() {
	log.Printf("bench-uploader started at %s\n", time.Now().String()[:19])
	if err := Command.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	log.Printf("bench-uploader ended at %s\n", time.Now().String()[:19])
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	kbts, err := ioutil.ReadFile(googleCloudStorageJSONKeyPath)
	if err != nil {
		return err
	}
	conf, err := google.JWTConfigFromJSON(
		kbts,
		storage.ScopeFullControl,
	)
	if err != nil {
		return err
	}
	ctx := context.Background()
	aclient, err := storage.NewAdminClient(ctx, googleCloudProjectName, cloud.WithTokenSource(conf.TokenSource(ctx)))
	if err != nil {
		return err
	}
	defer aclient.Close()

	if err := aclient.CreateBucket(context.Background(), googleCloudStorageBucketName, nil); err != nil {
		if !strings.Contains(err.Error(), "You already own this bucket. Please select another name") {
			return err
		}
	}

	sctx := context.Background()
	sclient, err := storage.NewClient(sctx, cloud.WithTokenSource(conf.TokenSource(sctx)))
	if err != nil {
		return err
	}
	defer sclient.Close()

	log.Printf("Uploading %s\n", csvResultPath)

	wc := sclient.Bucket(googleCloudStorageBucketName).Object(filepath.Base(csvResultPath)).NewWriter(context.Background())
	wc.ContentType = "text/plain"
	bts, err := ioutil.ReadFile(csvResultPath)
	if err != nil {
		return err
	}
	if _, err := wc.Write(bts); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}
