// Copyright © 2016 Shingo Omura <everpeace@gmail.com>
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

package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

//
// Flags
//
var bucketRegion string
var bucketName string

//
// commands
//
var assetFileNames = []string{
	"cluster.yml",
	"terraform.tfstate",
	"host_key",
	"host_key.pub",
	"worker_key",
	"worker_key.pub",
	"session_signing_key",
	"session_signing_key.pub",
	"authorized_worker_keys",
}

var save = &cobra.Command{
	Use:   "save",
	Short: "Save state files to S3",
	Long: `Put state files below to specified S3.  Keys to be stored in specified bucket are the same with filenames.
  - ` + strings.Join(assetFileNames[:], ", "),
	Run: RunSave,
}

var restore = &cobra.Command{
	Use:   "restore",
	Short: "Restore state files from S3",
	Long: `Restore state files below from specified S3. Keys pulled from specified bucket are the same with filenames.
  - ` + strings.Join(assetFileNames[:], ", "),
	Run: RunRestore,
}

func RunSave(cmd *cobra.Command, args []string) {
	if len(bucketName) < 1 || len(bucketRegion) < 1 {
		log.Panic("--bucket and --bucket-region are required.")
	}

	PutFilesToS3(bucketRegion, bucketName, cfgDir, assetFileNames)
}

func RunRestore(cmd *cobra.Command, args []string) {
	if len(bucketName) < 1 {
		log.Panic("--bucket is required.")
	}

	makeCfgDirIfNotExists()

	GetFilesFromS3(bucketRegion, bucketName, cfgDir, assetFileNames)
}

func init() {
	RootCmd.AddCommand(restore)
	restore.Flags().StringVar(&bucketName, "bucket", "", "S3 bucket name to which assets will be uploaded.")
	restore.Flags().StringVar(&bucketRegion, "bucket-region", "", "Region of S3 Bucket specified by --bucket")

	RootCmd.AddCommand(save)
	save.Flags().StringVar(&bucketName, "bucket", "", "S3 bucket name to which assets will be uploaded.")
	save.Flags().StringVar(&bucketRegion, "bucket-region", "", "Region of S3 Bucket specified by --bucket")
}
