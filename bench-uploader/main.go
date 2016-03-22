// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// bench-uploader uploads bench results to cloud storage.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/coreos/dbtester/remotestorage"
	"github.com/spf13/cobra"
)

// TODO: vendor gRPC and combine this into bench
// Currently, we get:
// panic: http: multiple registrations for /debug/requests

var (
	Command = &cobra.Command{
		Use:   "bench-uploader",
		Short: "Uploads to cloud storage.",
		RunE:  CommandFunc,
	}
)

var (
	from        string
	to          string
	isDirectory bool

	googleCloudProjectName        string
	googleCloudStorageJSONKeyPath string
	googleCloudStorageBucketName  string
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	Command.PersistentFlags().StringVar(&from, "from", "", "file to upload.")
	Command.PersistentFlags().StringVar(&to, "to", "", "file path to upload.")
	Command.PersistentFlags().BoolVar(&isDirectory, "directory", false, "'true' if uploading directory.")
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
	kbs, err := ioutil.ReadFile(googleCloudStorageJSONKeyPath)
	if err != nil {
		return err
	}
	u, err := remotestorage.NewGoogleCloudStorage(kbs, googleCloudProjectName)
	if err != nil {
		return err
	}
	if !isDirectory {
		if err := u.UploadFile(googleCloudStorageBucketName, from, to); err != nil {
			return err
		}
	} else {
		if err := u.UploadDir(googleCloudStorageBucketName, from, to); err != nil {
			return err
		}
	}
	return nil
}
