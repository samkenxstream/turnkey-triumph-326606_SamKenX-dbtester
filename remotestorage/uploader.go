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

package remotestorage

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

type Uploader interface {
	UploadFile(bucket, src, dst string, opts ...OpOption) error
	UploadDir(bucket, src, dst string, opts ...OpOption) error
}

// GoogleCloudStorage wraps Google Cloud Storage API.
type GoogleCloudStorage struct {
	JSONKey []byte
	Project string
	Config  *jwt.Config
}

func NewGoogleCloudStorage(key []byte, project string) (Uploader, error) {
	conf, err := google.JWTConfigFromJSON(
		key,
		storage.ScopeFullControl,
	)
	if err != nil {
		return nil, err
	}
	return &GoogleCloudStorage{
		JSONKey: key,
		Project: project,
		Config:  conf,
	}, nil
}

func (g *GoogleCloudStorage) UploadFile(bucket, src, dst string, opts ...OpOption) error {
	if g == nil {
		return fmt.Errorf("GoogleCloudStorage is nil")
	}
	ret := &Op{}
	ret.applyOpts(opts)

	ctx := context.Background()
	admin, err := storage.NewAdminClient(ctx, g.Project, cloud.WithTokenSource(g.Config.TokenSource(ctx)))
	if err != nil {
		return err
	}
	defer admin.Close()

	if err := admin.CreateBucket(context.Background(), bucket, nil); err != nil {
		if !strings.Contains(err.Error(), "You already own this bucket. Please select another name") {
			return err
		}
	}

	sctx := context.Background()
	client, err := storage.NewClient(sctx, cloud.WithTokenSource(g.Config.TokenSource(sctx)))
	if err != nil {
		return err
	}
	defer client.Close()

	wc := client.Bucket(bucket).Object(dst).NewWriter(context.Background())
	if ret.ContentType != "" {
		wc.ContentType = ret.ContentType
	}

	fmt.Println("UploadFile:")
	fmt.Println()
	fmt.Println(src)
	fmt.Println("--->")
	fmt.Println(dst)
	fmt.Println()

	bts, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile(%s) %v", src, err)
	}
	if _, err := wc.Write(bts); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("UploadFile Done!")
	fmt.Println()
	return nil
}

func (g *GoogleCloudStorage) UploadDir(bucket, src, dst string, opts ...OpOption) error {
	if g == nil {
		return fmt.Errorf("GoogleCloudStorage is nil")
	}
	ret := &Op{}
	ret.applyOpts(opts)

	ctx := context.Background()
	admin, err := storage.NewAdminClient(ctx, g.Project, cloud.WithTokenSource(g.Config.TokenSource(ctx)))
	if err != nil {
		return err
	}
	defer admin.Close()

	if err := admin.CreateBucket(context.Background(), bucket, nil); err != nil {
		if !strings.Contains(err.Error(), "You already own this bucket. Please select another name") {
			return err
		}
	}

	sctx := context.Background()
	client, err := storage.NewClient(sctx, cloud.WithTokenSource(g.Config.TokenSource(sctx)))
	if err != nil {
		return err
	}
	defer client.Close()

	fmap, err := walkRecursive(src)
	if err != nil {
		return err
	}

	fmt.Println("UploadDir:")
	donec, errc := make(chan struct{}), make(chan error)
	for source := range fmap {
		go func(source string) {
			s := strings.Replace(source, src, "", -1)
			tp := filepath.Join(dst, s)

			fmt.Println()
			fmt.Println(source)
			fmt.Println("--->")
			fmt.Println(tp)
			fmt.Println()

			wc := client.Bucket(bucket).Object(tp).NewWriter(context.Background())
			if ret.ContentType != "" {
				wc.ContentType = ret.ContentType
			}
			bts, err := ioutil.ReadFile(source)
			if err != nil {
				errc <- fmt.Errorf("ioutil.ReadFile(%s) %v", source, err)
				return
			}
			if _, err := wc.Write(bts); err != nil {
				errc <- err
				return
			}
			if err := wc.Close(); err != nil {
				errc <- err
				return
			}
			donec <- struct{}{}
		}(source)
	}

	cnt, num := 0, len(fmap)
	for cnt != num {
		select {
		case <-donec:
		case err := <-errc:
			return err
		}
		cnt++
	}
	fmt.Println()
	fmt.Println("UploadDir Done!")
	fmt.Println()
	return nil
}
