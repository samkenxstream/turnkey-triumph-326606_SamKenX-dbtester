// Copyright 2017 CoreOS, Inc.
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

package control

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/dbtester/pkg/remotestorage"
)

func step4UploadLogs(cfg Config) error {
	if err := uploadToGoogle(cfg.Log, cfg); err != nil {
		return err
	}
	if err := uploadToGoogle(cfg.DatasizeSummary, cfg); err != nil {
		return err
	}
	if err := uploadToGoogle(cfg.DataLatencyDistributionSummary, cfg); err != nil {
		return err
	}
	if err := uploadToGoogle(cfg.DataLatencyDistributionPercentile, cfg); err != nil {
		return err
	}
	if err := uploadToGoogle(cfg.DataLatencyDistributionAll, cfg); err != nil {
		return err
	}
	if err := uploadToGoogle(cfg.DataLatencyThroughputTimeseries, cfg); err != nil {
		return err
	}
	return nil
}

func uploadToGoogle(path string, cfg Config) error {
	if !exist(path) {
		return fmt.Errorf("%q does not exist", path)
	}
	u, err := remotestorage.NewGoogleCloudStorage([]byte(cfg.GoogleCloudStorageKey), cfg.GoogleCloudProjectName)
	if err != nil {
		return err
	}

	srcPath := path
	dstPath := filepath.Base(path)
	if !strings.HasPrefix(dstPath, cfg.TestName) {
		dstPath = fmt.Sprintf("%s-%s", cfg.TestName, dstPath)
	}
	dstPath = filepath.Join(cfg.GoogleCloudStorageSubDirectory, dstPath)

	var uerr error
	for k := 0; k < 30; k++ {
		if uerr = u.UploadFile(cfg.GoogleCloudStorageBucketName, srcPath, dstPath); uerr != nil {
			plog.Printf("#%d: error %v while uploading %q", k, uerr, path)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	return uerr
}
