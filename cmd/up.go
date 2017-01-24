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
	"strconv"
	"strings"

	"github.com/mumoshu/concourse-aws/concourse"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Spin up or Update concourse on aws",
	Long:  `This spins up scalable concourse ci servers interactively.  This command supports to update(implemented by terraform update) states of concourse ci servrers.`,
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	c, err := concourse.ConfigFromFile(prefixConfigDir("cluster.yml"))
	if err != nil && os.IsNotExist(err) {
		fmt.Printf("Creating cluster.yml")
		c = InteractivelyCreateConfig()
		WriteConfigFile(c, prefixConfigDir("cluster.yml"))
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
	prefix := AskForRequiredInput("Prefix", AskOptions{Default: "concourse-"})

	regions := ListRegions()
	region := AskForRequiredInput("Region", AskOptions{Candidates: regions, Validate: mustBeIncludedIn(regions), Default: "ap-northeast-1"})

	out, _ := exec.Command("./my-latest-ami.sh").CombinedOutput()
	latestAmiId := strings.TrimSpace(string(out))
	amiId := AskForRequiredInput("AMI ID", AskOptions{
		Default:    latestAmiId,
		Candidates: []string{latestAmiId},
	})

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

	accessibleCIDRS := AskForRequiredInput("AccessibleCIDRS(commma separated)", AskOptions{Default: fmt.Sprintf("%s/32", ObtainExternalIp())})

	dbInstanceClass := AskForRequiredInput("DB Instance Class", AskOptions{Default: "db.t2.micro"})

	webInstanceType := AskForRequiredInput("Concourse Web Instance Type", AskOptions{Default: "t2.micro"})
	workerInstanceType := AskForRequiredInput("Concourse Worker Instance Type", AskOptions{Default: "t2.micro"})

	asgMin := AskForRequiredInput("Min numbers of servers in ASG(Web, Worker)", AskOptions{Default: "0"})
	asgMax := AskForRequiredInput("Max numbers of servers in ASG(Web, Worker)", AskOptions{Default: "2"})
	webAsgDesired := AskForRequiredInput("Desired numbers of web servers in ASG", AskOptions{Default: "1"})
	workerAsgDesired := AskForRequiredInput("Desired numbers of servers in ASG", AskOptions{Default: "2"})

	possibleElbProtocols := []string{"http", "https"}
	defaultElbPorts := map[string]string{
		"http":  "80",
		"https": "443",
	}
	elbProtocol := AskForRequiredInput("Protocol for ELB", AskOptions{
		Default:    possibleElbProtocols[0],
		Candidates: possibleElbProtocols,
		Validate:   mustBeIncludedIn(possibleElbProtocols),
	})
	elbPort, err := strconv.Atoi(AskForRequiredInput("Port for ELB", AskOptions{
		Default: defaultElbPorts[elbProtocol],
	}))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	sslCertificateArn := ""
	if elbProtocol == "https" {
		sslCertificateArn = AskForRequiredInput("SSL ARN", AskOptions{Default: ""})
	}
	customExternalDomainName := AskForRequiredInput("Custom External Domain Name(just hit enter for skip, e.g. some.cool.com)", AskOptions{Default: ""})

	username := AskForRequiredInput("Basic Auth Username(just hit enter for skip)", AskOptions{Default: ""})
	password := ""
	if username != "" {
		password = AskForRequiredInput("Basic Auth Password", AskOptions{Default: "bar"})
	}

	ghClientId := AskForRequiredInput("Github Auth Client Id(just hit enter for skip)", AskOptions{Default: ""})
	ghClientSecret := ""
	ghOrgs := []string{}
	ghTeams := []string{}
	ghUsers := []string{}
	if ghClientId != "" {
		ghClientSecret = AskForRequiredInput("Github Auth Client Secret(just hit enter for skip)", AskOptions{Default: ""})
		ghOrgsInput := AskForRequiredInput("Github Auth Organizations(comma separated)", AskOptions{Default: ""})
		if ghOrgsInput != "" {
			ghOrgs = strings.Split(ghOrgsInput, ",")
		}

		ghTeamsInput := AskForRequiredInput("Github Auth Teams(comma separated, e.g. ORG/TEAM)", AskOptions{Default: ""})
		if ghTeamsInput != "" {
			ghTeams = strings.Split(ghTeamsInput, ",")
		}

		ghUsersInput := AskForRequiredInput("Github Auth Users(comma separated)", AskOptions{Default: ""})
		if ghUsersInput != "" {
			ghUsers = strings.Split(ghUsersInput, ",")
		}
	}

	if username == "" && ghClientId == "" {
		fmt.Println("WARNING WARNING WARNING WARNING WARNING")
		fmt.Println("!!!  No Authentication configured   !!!")
		fmt.Println("WARNING WARNING WARNING WARNING WARNING")
		possibleAnswers := []string{"y", "n"}
		yesOrNo := AskForRequiredInput("Do you really want to procceed?", AskOptions{Candidates: possibleAnswers, Validate: mustBeIncludedIn(possibleAnswers), Default: "n"})
		if yesOrNo == "n" {
			os.Exit(1)
		}
	}

	return &concourse.Config{
		Prefix:                   prefix,
		Region:                   region,
		KeyName:                  keyName,
		AccessibleCIDRS:          accessibleCIDRS,
		VpcId:                    vpcId,
		SubnetIds:                subnetIds,
		AvailabilityZones:        availabilityZones,
		DBInstanceClass:          dbInstanceClass,
		WebInstanceType:          webInstanceType,
		WorkerInstanceType:       workerInstanceType,
		AMI:                      amiId,
		AsgMin:                   asgMin,
		AsgMax:                   asgMax,
		WebAsgDesired:            webAsgDesired,
		WorkerAsgDesired:         workerAsgDesired,
		ElbProtocol:              elbProtocol,
		ElbPort:                  elbPort,
		CustomExternalDomainName: customExternalDomainName,
		SSLCertificateArn:        sslCertificateArn,
		BasicAuthUsername:        username,
		BasicAuthPassword:        password,
		GithubAuthClientId:       ghClientId,
		GithubAuthClientSecret:   ghClientSecret,
		GithubAuthOrganizations:  ghOrgs,
		GithubAuthTeams:          ghTeams,
		GithubAuthUsers:          ghUsers,
	}
}

func WriteConfigFile(config *concourse.Config, path string) {
	d, err := yaml.Marshal(&config)
	if err != nil {
		panic(err)
	}

	makeCfgDirIfNotExists()

	if ioutil.WriteFile(path, []byte(d), 0644) != nil {
		panic(err)
	}
}

func SSHGenKeyIfNotExist(keyFileName string) {
	if _, err := os.Stat(keyFileName); os.IsNotExist(err) {
		log.Println(fmt.Sprintf("generating ssh key: %s", keyFileName))
		args := []string{
			"-t", "rsa",
			"-f", keyFileName,
			"-N", "",
		}
		cmd := exec.Command("ssh-keygen", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
}

func TerraformRun(subcommand string, c *concourse.Config) {
	// auto ssh key creation
	SSHGenKeyIfNotExist(prefixConfigDir("host_key"))
	SSHGenKeyIfNotExist(prefixConfigDir("worker_key"))
	SSHGenKeyIfNotExist(prefixConfigDir("session_signing_key"))
	cp := exec.Command("cp",
		prefixConfigDir("worker_key.pub"),
		prefixConfigDir("authorized_worker_keys"))
	cp.Stdout = os.Stdout
	cp.Stderr = os.Stderr
	if err := cp.Run(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	useCustomExternalDomainName := 0
	if len(c.CustomExternalDomainName) > 0 {
		useCustomExternalDomainName = 1
	}
	useCustomElbPort := 0
	if !(c.ElbPort == 80 || c.ElbPort == 443) {
		useCustomElbPort = 1
	}

	// for backward compatibility
	// instance_type will be copied to web_instance_type and worker_instance_type only if they are not used.
	if c.InstanceType != "" {
		fmt.Println("WARNING: instance_type is deprecated. Use web_instance_type and worker_instance_type instead")
		if c.WebInstanceType == "" {
			c.WebInstanceType = c.InstanceType
		}
		if c.WorkerInstanceType == "" {
			c.WorkerInstanceType = c.InstanceType
		}
	}

	args := []string{
		subcommand,
		"-state", fmt.Sprintf("%s%s", cfgDir, "terraform.tfstate"),
		"-var", fmt.Sprintf("aws_region=%s", c.Region),
		"-var", fmt.Sprintf("availability_zones=%s", strings.Join(c.AvailabilityZones, ",")),
		"-var", fmt.Sprintf("key_name=%s", c.KeyName),
		"-var", fmt.Sprintf("subnet_id=%s", strings.Join(c.SubnetIds, ",")),
		"-var", fmt.Sprintf("vpc_id=%s", c.VpcId),
		"-var", fmt.Sprintf("db_instance_class=%s", c.DBInstanceClass),
		"-var", fmt.Sprintf("web_instance_type=%s", c.WebInstanceType),
		"-var", fmt.Sprintf("worker_instance_type=%s", c.WorkerInstanceType),
		"-var", "db_username=concourse",
		"-var", "db_password=concourse",
		"-var", fmt.Sprintf("db_subnet_ids=%s", strings.Join(c.SubnetIds, ",")),
		"-var", fmt.Sprintf("tsa_host_key=%s", prefixConfigDir("host_key")),
		"-var", fmt.Sprintf("session_signing_key=%s", prefixConfigDir("session_signing_key")),
		"-var", fmt.Sprintf("tsa_authorized_keys=%s", prefixConfigDir("worker_key.pub")),
		"-var", fmt.Sprintf("tsa_public_key=%s", prefixConfigDir("host_key.pub")),
		"-var", fmt.Sprintf("tsa_worker_private_key=%s", prefixConfigDir("worker_key")),
		"-var", fmt.Sprintf("ami=%s", c.AMI),
		"-var", fmt.Sprintf("in_access_allowed_cidrs=%s", c.AccessibleCIDRS),
		"-var", fmt.Sprintf("elb_listener_lb_protocol=%s", c.ElbProtocol),
		"-var", fmt.Sprintf("elb_listener_lb_port=%d", c.ElbPort),
		"-var", fmt.Sprintf("use_custom_elb_port=%d", useCustomElbPort),
		"-var", fmt.Sprintf("ssl_certificate_arn=%s", c.SSLCertificateArn),
		"-var", fmt.Sprintf("use_custom_external_domain_name=%d", useCustomExternalDomainName),
		"-var", fmt.Sprintf("custom_external_domain_name=%s", c.CustomExternalDomainName),
		"-var", fmt.Sprintf("worker_instance_profile=%s", c.WorkerInstanceProfile),
		"-var", fmt.Sprintf("basic_auth_username=%s", c.BasicAuthUsername),
		"-var", fmt.Sprintf("basic_auth_password=%s", c.BasicAuthPassword),
		"-var", fmt.Sprintf("github_auth_client_id=%s", c.GithubAuthClientId),
		"-var", fmt.Sprintf("github_auth_client_secret=%s", c.GithubAuthClientSecret),
	}

	if len(c.Prefix) > 0 {
		args = append(args,
			"-var", fmt.Sprintf("prefix=%s", c.Prefix),
		)
	}

	if len(c.AsgMin) > 0 {
		args = append(args,
			"-var", fmt.Sprintf("asg_min=%s", c.AsgMin),
		)
	}
	if len(c.AsgMax) > 0 {
		args = append(args,
			"-var", fmt.Sprintf("asg_max=%s", c.AsgMax),
		)
	}
	if len(c.WebAsgDesired) > 0 {
		args = append(args,
			"-var", fmt.Sprintf("web_asg_desired=%s", c.WebAsgDesired),
		)
	}
	if len(c.WorkerAsgDesired) > 0 {
		args = append(args,
			"-var", fmt.Sprintf("worker_asg_desired=%s", c.WorkerAsgDesired),
		)
	}

	if len(c.GithubAuthOrganizations) > 0 {
		args = append(args,
			"-var", fmt.Sprintf("github_auth_organizations=%s", strings.Join(c.GithubAuthOrganizations, ",")),
		)
	}
	if len(c.GithubAuthTeams) > 0 {
		args = append(args,
			"-var", fmt.Sprintf("github_auth_teams=%s", strings.Join(c.GithubAuthTeams, ",")),
		)
	}
	if len(c.GithubAuthUsers) > 0 {
		args = append(args,
			"-var", fmt.Sprintf("github_auth_users=%s", strings.Join(c.GithubAuthUsers, ",")),
		)
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
