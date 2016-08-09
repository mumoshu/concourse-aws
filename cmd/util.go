package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func PutFilesToS3(region string, bucketName string, path string, filenames []string) {
	svc := s3.New(session.New(), &aws.Config{Region: aws.String(region)})

	for _, filename := range filenames {
		file, err := os.Open(fmt.Sprintf("%s%s", path, filename))
		if err != nil {
			log.Panic(err.Error())
		}
		defer file.Close()

		resp, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(filename),
			Body:   file,
		})

		if err != nil {
			log.Panic(err.Error())
		}

		log.Println(awsutil.StringValue(resp))
	}
}

func GetFilesFromS3(region string, bucketName string, path string, filenames []string) {
	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(region)}))

	for _, filename := range filenames {
		file, err := os.Create(fmt.Sprintf("%s%s", path, filename))
		if err != nil {
			log.Panic("Failed to create file", err)
		}
		defer file.Close()

		numBytes, err := downloader.Download(file, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(filename),
		})

		if err != nil {
			log.Panic(err.Error())
		}

		log.Println("Downloaded file: ", file.Name(), " (", numBytes, "bytes)")
	}
}
