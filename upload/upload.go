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

package upload

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/coreos/dbtester/remotestorage"
	"github.com/spf13/cobra"
)

// TODO: vendor gRPC and combine this into bench
// Currently, we get:
// panic: http: multiple registrations for /debug/requests

var (
	Command = &cobra.Command{
		Use:   "upload",
		Short: "Uploads to cloud storage.",
		RunE:  CommandFunc,
	}
)

var (
	from        string
	to          string
	isDirectory bool

	googleCloudProjectName string
	keyPath                string
	bucket                 string
)

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	Command.PersistentFlags().StringVar(&from, "from", "", "file to upload.")
	Command.PersistentFlags().StringVar(&to, "to", "", "file path to upload.")
	Command.PersistentFlags().BoolVar(&isDirectory, "directory", false, "'true' if uploading directory.")
	Command.PersistentFlags().StringVar(&googleCloudProjectName, "google-cloud-project-name", "", "Google cloud project name.")
	Command.PersistentFlags().StringVar(&keyPath, "key-path", "", "Path of key file.")
	Command.PersistentFlags().StringVar(&bucket, "bucket", "", "Bucket name in cloud storage.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	log.Println("opening key", keyPath)
	kbs, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("error when opening key %s(%v)", keyPath, err)
	}
	u, err := remotestorage.NewGoogleCloudStorage(kbs, googleCloudProjectName)
	if err != nil {
		return fmt.Errorf("error when NewGoogleCloudStorage %s(%v)", googleCloudProjectName, err)
	}
	if !isDirectory {
		if err := u.UploadFile(bucket, from, to); err != nil {
			return err
		}
	} else {
		if err := u.UploadDir(bucket, from, to); err != nil {
			return err
		}
	}
	return nil
}
