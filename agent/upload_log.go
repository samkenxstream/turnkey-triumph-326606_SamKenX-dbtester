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

package agent

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/etcd-io/dbtester/dbtesterpb"
	"github.com/etcd-io/dbtester/pkg/remotestorage"

	"go.uber.org/zap"
)

// uploadLog starts cetcd. This assumes that etcd is already started.
func uploadLog(fs *flags, t *transporterServer) error {
	t.lg.Info(
		"stopped collecting metrics, now uploading logs to storage",
		zap.String("gcp-project-name", t.req.ConfigClientMachineInitial.GoogleCloudProjectName),
	)
	u, err := remotestorage.NewGoogleCloudStorage(t.lg, []byte(t.req.ConfigClientMachineInitial.GoogleCloudStorageKey), t.req.ConfigClientMachineInitial.GoogleCloudProjectName)
	if err != nil {
		return err
	}

	var uerr error

	{
		srcDatabaseLogPath := fs.databaseLog
		dstDatabaseLogPath := filepath.Base(fs.databaseLog)
		if !strings.HasPrefix(filepath.Base(fs.databaseLog), t.req.DatabaseTag) {
			dstDatabaseLogPath = fmt.Sprintf("%s-%d-%s", t.req.DatabaseTag, t.req.IPIndex+1, filepath.Base(fs.databaseLog))
		}
		dstDatabaseLogPath = filepath.Join(t.req.ConfigClientMachineInitial.GoogleCloudStorageSubDirectory, dstDatabaseLogPath)
		t.lg.Info("uploading database log", zap.String("source", srcDatabaseLogPath), zap.String("destination", dstDatabaseLogPath))
		for k := 0; k < 30; k++ {
			if uerr = u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcDatabaseLogPath, dstDatabaseLogPath); uerr != nil {
				t.lg.Warn("upload error; retrying...", zap.Error(uerr))
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}
		if uerr != nil {
			return uerr
		}
	}

	{
		if t.req.DatabaseID == dbtesterpb.DatabaseID_zetcd__beta ||
			t.req.DatabaseID == dbtesterpb.DatabaseID_cetcd__beta {
			dpath := fs.databaseLog + "-" + t.req.DatabaseID.String()
			srcDatabaseLogPath2 := dpath
			dstDatabaseLogPath2 := filepath.Base(dpath)
			if !strings.HasPrefix(filepath.Base(dpath), t.req.DatabaseTag) {
				dstDatabaseLogPath2 = fmt.Sprintf("%s-%d-%s", t.req.DatabaseTag, t.req.IPIndex+1, filepath.Base(dpath))
			}
			dstDatabaseLogPath2 = filepath.Join(t.req.ConfigClientMachineInitial.GoogleCloudStorageSubDirectory, dstDatabaseLogPath2)
			t.lg.Info("uploading proxy database log", zap.String("source", srcDatabaseLogPath2), zap.String("destination", dstDatabaseLogPath2))
			for k := 0; k < 30; k++ {
				if uerr = u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcDatabaseLogPath2, dstDatabaseLogPath2); uerr != nil {
					t.lg.Warn("upload error; retrying...", zap.Error(uerr))
					time.Sleep(2 * time.Second)
					continue
				}
				break
			}
			if uerr != nil {
				return uerr
			}
		}
	}

	{
		srcSysMetricsDataPath := fs.systemMetricsCSV
		dstSysMetricsDataPath := filepath.Base(fs.systemMetricsCSV)
		if !strings.HasPrefix(filepath.Base(fs.systemMetricsCSV), t.req.DatabaseTag) {
			dstSysMetricsDataPath = fmt.Sprintf("%s-%d-%s", t.req.DatabaseTag, t.req.IPIndex+1, filepath.Base(fs.systemMetricsCSV))
		}
		dstSysMetricsDataPath = filepath.Join(t.req.ConfigClientMachineInitial.GoogleCloudStorageSubDirectory, dstSysMetricsDataPath)
		t.lg.Info("uploading system metrics", zap.String("source", srcSysMetricsDataPath), zap.String("destination", dstSysMetricsDataPath))
		for k := 0; k < 30; k++ {
			if uerr := u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcSysMetricsDataPath, dstSysMetricsDataPath); uerr != nil {
				t.lg.Warn("upload error; retrying...", zap.Error(uerr))
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}
		if uerr != nil {
			return uerr
		}
	}

	{
		srcSysMetricsInterpolatedDataPath := fs.systemMetricsCSVInterpolated
		dstSysMetricsInterpolatedDataPath := filepath.Base(fs.systemMetricsCSVInterpolated)
		if !strings.HasPrefix(filepath.Base(fs.systemMetricsCSVInterpolated), t.req.DatabaseTag) {
			dstSysMetricsInterpolatedDataPath = fmt.Sprintf("%s-%d-%s", t.req.DatabaseTag, t.req.IPIndex+1, filepath.Base(fs.systemMetricsCSVInterpolated))
		}
		dstSysMetricsInterpolatedDataPath = filepath.Join(t.req.ConfigClientMachineInitial.GoogleCloudStorageSubDirectory, dstSysMetricsInterpolatedDataPath)
		t.lg.Info("uploading interpolated system metrics", zap.String("source", srcSysMetricsInterpolatedDataPath), zap.String("destination", dstSysMetricsInterpolatedDataPath))
		for k := 0; k < 30; k++ {
			if uerr := u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcSysMetricsInterpolatedDataPath, dstSysMetricsInterpolatedDataPath); uerr != nil {
				t.lg.Warn("upload error; retrying...", zap.Error(uerr))
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}
		if uerr != nil {
			return uerr
		}
	}

	{
		srcAgentLogPath := fs.agentLog
		dstAgentLogPath := filepath.Base(fs.agentLog)
		if !strings.HasPrefix(filepath.Base(fs.agentLog), t.req.DatabaseTag) {
			dstAgentLogPath = fmt.Sprintf("%s-%d-%s", t.req.DatabaseTag, t.req.IPIndex+1, filepath.Base(fs.agentLog))
		}
		dstAgentLogPath = filepath.Join(t.req.ConfigClientMachineInitial.GoogleCloudStorageSubDirectory, dstAgentLogPath)
		t.lg.Info("uploading agent log", zap.String("source", srcAgentLogPath), zap.String("destination", dstAgentLogPath))
		for k := 0; k < 30; k++ {
			if uerr := u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcAgentLogPath, dstAgentLogPath); uerr != nil {
				t.lg.Warn("upload error; retrying...", zap.Error(uerr))
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}
	}

	return uerr
}
