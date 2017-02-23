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

	"github.com/coreos/dbtester/dbtesterpb"
	"github.com/coreos/dbtester/pkg/remotestorage"
)

// uploadLog starts cetcd. This assumes that etcd is already started.
func uploadLog(fs *flags, t *transporterServer) error {
	plog.Infof("stopped collecting metrics; uploading logs to storage %q", t.req.ConfigClientMachineInitial.GoogleCloudProjectName)
	u, err := remotestorage.NewGoogleCloudStorage([]byte(t.req.ConfigClientMachineInitial.GoogleCloudStorageKey), t.req.ConfigClientMachineInitial.GoogleCloudProjectName)
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
		plog.Infof("uploading database log [%q -> %q]", srcDatabaseLogPath, dstDatabaseLogPath)
		for k := 0; k < 30; k++ {
			if uerr = u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcDatabaseLogPath, dstDatabaseLogPath); uerr != nil {
				plog.Warningf("UploadFile error... sleep and retry... (%v)", uerr)
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
		}
		if uerr != nil {
			return uerr
		}
	}

	{
		if t.req.DatabaseID == dbtesterpb.DatabaseID_zetcd__beta || t.req.DatabaseID == dbtesterpb.DatabaseID_cetcd__beta {
			dpath := fs.databaseLog + "-" + t.req.DatabaseID.String()
			srcDatabaseLogPath2 := dpath
			dstDatabaseLogPath2 := filepath.Base(dpath)
			if !strings.HasPrefix(filepath.Base(dpath), t.req.DatabaseTag) {
				dstDatabaseLogPath2 = fmt.Sprintf("%s-%d-%s", t.req.DatabaseTag, t.req.IPIndex+1, filepath.Base(dpath))
			}
			dstDatabaseLogPath2 = filepath.Join(t.req.ConfigClientMachineInitial.GoogleCloudStorageSubDirectory, dstDatabaseLogPath2)
			plog.Infof("uploading proxy-database log [%q -> %q]", srcDatabaseLogPath2, dstDatabaseLogPath2)
			for k := 0; k < 30; k++ {
				if uerr = u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcDatabaseLogPath2, dstDatabaseLogPath2); uerr != nil {
					plog.Warningf("UploadFile error... sleep and retry... (%v)", uerr)
					time.Sleep(2 * time.Second)
					continue
				} else {
					break
				}
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
		plog.Infof("uploading system metrics data [%q -> %q]", srcSysMetricsDataPath, dstSysMetricsDataPath)
		for k := 0; k < 30; k++ {
			if uerr := u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcSysMetricsDataPath, dstSysMetricsDataPath); uerr != nil {
				plog.Warningf("upload error... sleep and retry... (%v)", uerr)
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
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
		plog.Infof("uploading system metrics interpolated data [%q -> %q]", srcSysMetricsInterpolatedDataPath, dstSysMetricsInterpolatedDataPath)
		for k := 0; k < 30; k++ {
			if uerr := u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcSysMetricsInterpolatedDataPath, dstSysMetricsInterpolatedDataPath); uerr != nil {
				plog.Warningf("upload error... sleep and retry... (%v)", uerr)
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
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
		plog.Infof("uploading agent logs [%q -> %q]", srcAgentLogPath, dstAgentLogPath)
		for k := 0; k < 30; k++ {
			if uerr := u.UploadFile(t.req.ConfigClientMachineInitial.GoogleCloudStorageBucketName, srcAgentLogPath, dstAgentLogPath); uerr != nil {
				plog.Warningf("UploadFile error... sleep and retry... (%v)", uerr)
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
		}
	}

	return uerr
}
