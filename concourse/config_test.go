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
`,
		expectedConfig: Config{
			Region:          "ap-northeast-1",
			KeyName:         "cw_kuoka",
			SubnetIds:       []string{"subnet-11111914", "subnet-2222fc48"},
			AccessibleCIDRS: "123.123.234.234/32,234.234.234.234/32",
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
