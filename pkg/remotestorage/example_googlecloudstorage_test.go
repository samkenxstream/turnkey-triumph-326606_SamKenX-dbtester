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

package remotestorage_test

import (
	"io/ioutil"
	"log"

	"github.com/coreos/dbtester/pkg/remotestorage"
)

func ExampleGoogleCloudStorage_UploadFile() {
	kbs, err := ioutil.ReadFile("key.json")
	if err != nil {
		log.Fatal(err)
	}
	u, err := remotestorage.NewGoogleCloudStorage(kbs, "my-project")
	if err != nil {
		log.Fatal(err)
	}
	if err := u.UploadFile("test-bucket", "agent.log", "dir/agent.log", remotestorage.WithContentType("text/plain")); err != nil {
		log.Fatal(err)
	}
}

func ExampleGoogleCloudStorage_UploadDir() {
	kbs, err := ioutil.ReadFile("key.json")
	if err != nil {
		log.Fatal(err)
	}
	u, err := remotestorage.NewGoogleCloudStorage(kbs, "my-project")
	if err != nil {
		log.Fatal(err)
	}
	if err := u.UploadDir("test-bucket", "articles", "articles"); err != nil {
		log.Fatal(err)
	}
}
