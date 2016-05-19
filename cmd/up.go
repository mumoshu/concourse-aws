// Copyright © 2016 Yusuke KUOKA <ykuoka@gmail.com>
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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mumoshu/concourse-aws/concourse"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: Run,
}

func Run(cmd *cobra.Command, args []string) {
	c, err := concourse.ConfigFromFile("cluster.yml")
	if err != nil && os.IsNotExist(err) {
		fmt.Printf("Creating cluster.yml")
		c = InteractivelyCreateConfig()
		WriteConfigFile(c, "cluster.yml")
	}
	//	fmt.Printf("config:%+v", c)
	TerraformRun("plan", c)
	possibleAnswers := []string{"y", "n"}
	yesOrNo := AskForRequiredInput("Apply?", AskOptions{Candidates: possibleAnswers, Validate: mustBeIncludedIn(possibleAnswers), Default: "y"})
	if yesOrNo == "y" {
		TerraformRun("apply", c)
	}
}

func mustBeIncludedIn(candidates []string) func(string) error {
	return func(item string) error {
		for _, r := range candidates {
			if r == item {
				return nil
			}
		}
		return fmt.Errorf("%s must be one of: %v", item, candidates)
	}
}

func InteractivelyCreateConfig() *concourse.Config {
	regions := ListRegions()
	region := AskForRequiredInput("Region", AskOptions{Candidates: regions, Validate: mustBeIncludedIn(regions), Default: "ap-northeast-1"})

	keys := ListKeys(region)
	keyName := AskForRequiredInput("KeyName", AskOptions{Candidates: keys, Validate: mustBeIncludedIn(keys)})

	vpcIds := ListVPCs(region)
	vpcId := AskForRequiredInput("VPC ID", AskOptions{Candidates: vpcIds, Validate: mustBeIncludedIn(vpcIds)})

	candidateZones := ListAvailabilityZones(region)
	subnetIds := []string{}
	availabilityZones := []string{}

	for i := 1; i <= 2; i++ {
		az := AskForRequiredInput(fmt.Sprintf("Availability Zone %d", i), AskOptions{Candidates: candidateZones, Validate: mustBeIncludedIn(candidateZones)})
		b := candidateZones[:0]
		for _, x := range candidateZones {
			if x != az {
				b = append(b, x)
			}
		}

		candidateSubnets := ListSubnets(region, vpcId, az)
		subnetId := AskForRequiredInput(fmt.Sprintf("Subnet %d", i), AskOptions{Candidates: candidateSubnets, Validate: mustBeIncludedIn(candidateSubnets)})

		availabilityZones = append(availabilityZones, az)
		subnetIds = append(subnetIds, subnetId)
	}

	accessibleCIDR := AskForRequiredInput("AccessibleCIDR", AskOptions{Default: fmt.Sprintf("%s/32", ObtainExternalIp())})

	dbInstanceClass := AskForRequiredInput("DB Instance Class", AskOptions{Default: "db.t2.micro"})
	instanceType := AskForRequiredInput("Concourse Instance Type", AskOptions{Default: "t2.micro"})

	amiId := AskForRequiredInput("AMI ID", AskOptions{})

	username := AskForRequiredInput("Basic Auth Username", AskOptions{Default: "foo"})
	password := AskForRequiredInput("Basic Auth Password", AskOptions{Default: "bar"})

	return &concourse.Config{
		Region:            region,
		KeyName:           keyName,
		AccessibleCIDR:    accessibleCIDR,
		VpcId:             vpcId,
		SubnetIds:         subnetIds,
		AvailabilityZones: availabilityZones,
		DBInstanceClass:   dbInstanceClass,
		InstanceType:      instanceType,
		AMI:               amiId,
		BasicAuthUsername: username,
		BasicAuthPassword: password,
	}
}

func WriteConfigFile(config *concourse.Config, path string) {
	d, err := yaml.Marshal(&config)
	if err != nil {
		panic(err)
	}

	if ioutil.WriteFile(path, []byte(d), 0644) != nil {
		panic(err)
	}
}

func TerraformRun(subcommand string, c *concourse.Config) {
	args := []string{
		subcommand,
		"-var", fmt.Sprintf("aws_region=%s", c.Region),
		"-var", fmt.Sprintf("availability_zones=%s", strings.Join(c.AvailabilityZones, ",")),
		"-var", fmt.Sprintf("key_name=%s", c.KeyName),
		"-var", fmt.Sprintf("subnet_id=%s", strings.Join(c.SubnetIds, ",")),
		"-var", fmt.Sprintf("vpc_id=%s", c.VpcId),
		"-var", fmt.Sprintf("db_instance_class=%s", c.DBInstanceClass),
		"-var", fmt.Sprintf("instance_type=%s", c.InstanceType),
		"-var", "db_username=concourse",
		"-var", "db_password=concourse",
		"-var", fmt.Sprintf("db_subnet_ids=%s", strings.Join(c.SubnetIds, ",")),
		"-var", "tsa_host_key=host_key",
		"-var", "session_signing_key=session_signing_key",
		"-var", "tsa_authorized_keys=worker_key.pub",
		"-var", "tsa_public_key=host_key.pub",
		"-var", "tsa_worker_private_key=worker_key",
		"-var", fmt.Sprintf("ami=%s", c.AMI),
		"-var", fmt.Sprintf("in_access_allowed_cidr=%s", c.AccessibleCIDR),
		"-var", fmt.Sprintf("worker_instance_profile=%s", c.WorkerInstanceProfile),
		"-var", fmt.Sprintf("basic_auth_username=%s", c.BasicAuthUsername),
		"-var", fmt.Sprintf("basic_auth_password=%s", c.BasicAuthPassword),
	}
	log.Println("Running terraform get")
	get := exec.Command("terraform", "get")
	get.Stdout = os.Stdout
	get.Stderr = os.Stderr
	getErr := get.Run()
	if getErr != nil {
		log.Fatal(getErr)
		panic(getErr)
	}

	log.Println(fmt.Sprintf("Running terraform with args: %s", args))
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func init() {
	RootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
