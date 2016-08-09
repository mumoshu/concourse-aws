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

// putStatesToS3 represents the up command
var putStatesToS3 = &cobra.Command{
	Use:   "put-states-to-s3",
	Short: "Put state files to S3",
	Long: `Put state files below to specified S3

	- ` + strings.Join(assetFileNames[:], ","),
	Run: RunPutStatesToS3,
}

// getStatesFromS3 represents the up command
var getStatesFromS3 = &cobra.Command{
	Use:   "get-states-from-s3",
	Short: "Get state files from S3",
	Long: `Get state files below from specified S3

	- ` + strings.Join(assetFileNames[:], ","),
	Run: RunGetStatesFromS3,
}

func RunPutStatesToS3(cmd *cobra.Command, args []string) {
	if len(bucketName) < 1 || len(bucketRegion) < 1 {
		log.Panic("--bucket and --bucket-region are required.")
	}

	PutFilesToS3(bucketRegion, bucketName, cfgDir, assetFileNames)
}

func RunGetStatesFromS3(cmd *cobra.Command, args []string) {
	if len(bucketName) < 1 {
		log.Panic("--bucket is required.")
	}

	makeCfgDirIfNotExists()

	GetFilesFromS3(bucketRegion, bucketName, cfgDir, assetFileNames)
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// putStatesToS3.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// putStatesToS3.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.AddCommand(putStatesToS3)
	putStatesToS3.Flags().StringVar(&bucketName, "bucket", "", "S3 bucket name to which assets will be uploaded.")
	putStatesToS3.Flags().StringVar(&bucketRegion, "bucket-region", "", "Region of S3 Bucket specified by --bucket")

	RootCmd.AddCommand(getStatesFromS3)
	getStatesFromS3.Flags().StringVar(&bucketName, "bucket", "", "S3 bucket name to which assets will be uploaded.")
	getStatesFromS3.Flags().StringVar(&bucketRegion, "bucket-region", "", "Region of S3 Bucket specified by --bucket")
}
