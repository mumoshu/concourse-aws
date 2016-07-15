package concourse

import (
	"reflect"
	"testing"
)

var validConfigs = []struct {
	providedYaml   string
	expectedConfig Config
}{
	{
		providedYaml: `
---
region: ap-northeast-1
key_name: cw_kuoka
subnet_ids:
  - subnet-11111914
  - subnet-2222fc48
accessible_cidrs: 123.123.234.234/32,234.234.234.234/32
github_auth_client_id: dummydummy
github_auth_client_secret: dummydummydummy
github_auth_organizations: [org1, org2]
github_auth_teams: [org3/team1, org3/team2]
github_auth_users: []
elb_protocol: "https"
elb_port: 443
custom_external_domain_name: "some.where"
ssl_certificate_arn: "arn://dummydummy"
`,
		expectedConfig: Config{
			Region:                   "ap-northeast-1",
			KeyName:                  "cw_kuoka",
			SubnetIds:                []string{"subnet-11111914", "subnet-2222fc48"},
			AccessibleCIDRS:          "123.123.234.234/32,234.234.234.234/32",
			ElbProtocol:              "https",
			ElbPort:                  443,
			CustomExternalDomainName: "some.where",
			SSLCertificateArn:        "arn://dummydummy",
			GithubAuthClientId:       "dummydummy",
			GithubAuthClientSecret:   "dummydummydummy",
			GithubAuthOrganizations:  []string{"org1", "org2"},
			GithubAuthTeams:          []string{"org3/team1", "org3/team2"},
			GithubAuthUsers:          []string{},
		},
	},
}

func TestParsing(t *testing.T) {
	for _, c := range validConfigs {
		actual, err := ConfigFromString(c.providedYaml)
		if err != nil {
			t.Errorf("Failed to parse config: %s: %v", c.providedYaml, err)
		}
		if !reflect.DeepEqual(actual, c.expectedConfig) {
			t.Errorf("Produced config does not match against expected config. Produced:\n%+v\nExpected:\n%+v\n", actual, c.expectedConfig)
		}
	}
}
