package concourse

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Region                string
	KeyName               string   `yaml:"key_name"`
	SubnetIds             []string `yaml:"subnet_ids"`
	VpcId                 string   `yaml:"vpc_id"`
	AvailabilityZones     []string `yaml:"availability_zones"`
	AccessibleCIDR        string   `yaml:"accessible_cidr"`
	DBInstanceClass       string   `yaml:"db_instance_class"`
	InstanceType          string   `yaml:"instance_type"`
	WorkerInstanceProfile string   `yaml:"worker_instance_profile"`
	AMI                   string   `yaml:"ami_id"`
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
