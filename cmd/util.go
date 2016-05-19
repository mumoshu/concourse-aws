package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func ListRegions() []string {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})

	params := &ec2.DescribeRegionsInput{
		DryRun: aws.Bool(false),
	}
	resp, err := svc.DescribeRegions(params)

	regions := []string{}
	for _, r := range resp.Regions {
		regions = append(regions, *r.RegionName)
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return regions
}

func ListKeys(region string) []string {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &ec2.DescribeKeyPairsInput{
		DryRun: aws.Bool(false),
	}
	resp, err := svc.DescribeKeyPairs(params)

	keyPairs := []string{}
	for _, k := range resp.KeyPairs {
		keyPairs = append(keyPairs, *k.KeyName)
	}

	if err != nil {
		fmt.Println(err.Error())
		return []string{}
	}
	return keyPairs
}

func ListVPCs(region string) []string {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &ec2.DescribeVpcsInput{
		DryRun: aws.Bool(false),
	}
	resp, err := svc.DescribeVpcs(params)

	if err != nil {
		panic(err)
	}

	vpcs := []string{}
	for _, c := range resp.Vpcs {
		vpcs = append(vpcs, *c.VpcId)
	}

	return vpcs

}

func ListAvailabilityZones(region string) []string {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &ec2.DescribeAvailabilityZonesInput{
		DryRun: aws.Bool(false),
	}
	resp, err := svc.DescribeAvailabilityZones(params)

	if err != nil {
		panic(err)
	}

	azs := []string{}
	for _, z := range resp.AvailabilityZones {
		azs = append(azs, *z.ZoneName)
	}

	return azs

}

func ListSubnets(region string, vpcId string, az string) []string {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &ec2.DescribeSubnetsInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("availabilityZone"),
				Values: []*string{
					aws.String(az),
				},
			},
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcId),
				},
			},
		},
	}
	resp, err := svc.DescribeSubnets(params)

	if err != nil {
		panic(err)
	}

	subnets := []string{}
	for _, s := range resp.Subnets {
		subnets = append(subnets, *s.SubnetId)
	}

	return subnets

}
