package concourse

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Prefix            string   `yaml:"prefix"`
	Region            string   `yaml:"region"`
	KeyName           string   `yaml:"key_name"`
	SubnetIds         []string `yaml:"subnet_ids"`
	VpcId             string   `yaml:"vpc_id"`
	AvailabilityZones []string `yaml:"availability_zones"`
	AccessibleCIDRS   string   `yaml:"accessible_cidrs"`
	DBInstanceClass   string   `yaml:"db_instance_class"`
	DBEngineVersion   string   `yaml:"db_engine_version"`
	// Deprecated: Use WebInstanceType and WorkerInstanceType instead.
	InstanceType             string   `yaml:"instance_type"`
	WebInstanceType          string   `yaml:"web_instance_type"`
	WorkerInstanceType       string   `yaml:"worker_instance_type"`
	WorkerInstanceProfile    string   `yaml:"worker_instance_profile"`
	AMI                      string   `yaml:"ami_id"`
	AsgMin                   string   `yaml:"asg_min"`
	AsgMax                   string   `yaml:"asg_max"`
	WebAsgDesired            string   `yaml:"web_asg_desired"`
	WorkerAsgDesired         string   `yaml:"worker_asg_desired"`
	ElbProtocol              string   `yaml:"elb_protocol"`
	ElbPort                  int      `yaml:"elb_port"`
	CustomExternalDomainName string   `yaml:"custom_external_domain_name"`
	SSLCertificateArn        string   `yaml:"ssl_certificate_arn"`
	BasicAuthUsername        string   `yaml:"basic_auth_username"`
	BasicAuthPassword        string   `yaml:"basic_auth_password"`
	GithubAuthClientId       string   `yaml:"github_auth_client_id"`
	GithubAuthClientSecret   string   `yaml:"github_auth_client_secret"`
	GithubAuthOrganizations  []string `yaml:"github_auth_organizations"`
	GithubAuthTeams          []string `yaml:"github_auth_teams"`
	GithubAuthUsers          []string `yaml:"github_auth_users"`
}

func ConfigFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c, err := ConfigFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("file %s: %v", filename, err)
	}

	return &c, nil
}

func ConfigFromString(data string) (Config, error) {
	return ConfigFromBytes([]byte(data))
}

func ConfigFromBytes(data []byte) (Config, error) {
	c := Config{}

	err := yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return c, err
}
