package main

import (
	//"errors"
	"errors"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AWSDefault : a struct to save default value
type AWSDefault struct {
	CredentialPath string
	Profile        string
	Region         string
	Bucket         string
	Prefix         string
}

// GetCredentialChain : return a credential from Environment , ec2role or shared credential. If uses shared credential,
// provides credentialPath and Profile in parameters
func GetCredentialChain(credentialPath, profile string) *aws.Config {
	config := aws.NewConfig()
	ec2m := ec2metadata.New(session.New(), config)
	var ProviderList []credentials.Provider = []credentials.Provider{
		&ec2rolecreds.EC2RoleProvider{
			Client: ec2m,
		},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{Filename: credentialPath, Profile: profile},
	}
	creds := credentials.NewChainCredentials(ProviderList)
	config.WithCredentials(creds)
	return config
}

//NewS3Service : Create a new session of S3 service
func NewS3Service(credentialPath, profile, region string) (*s3.S3, error) {
	config := GetCredentialChain(credentialPath, profile)
	config.WithRegion(region)
	s3Service := s3.New(session.New(), config)
	if s3Service != nil {
		return s3Service, nil
	}
	return nil, errors.New("s3 Service create fail")

}

/********
AWS Operations
**********/
// S3Request : A collection of S3 Request Functions
func S3Request(service *s3.S3, input interface{}) (output interface{}, err error) {
	if nil == service {
		return nil, errors.New("s3 service is not available")
	}
	itype := reflect.TypeOf(input)
	switch itype.Name() {
	case "PutObjectInput":
		p := input.(s3.PutObjectInput)
		output, err = service.PutObject(&p)
	case "GetObjectInput":
		p := input.(s3.GetObjectInput)
		output, err = service.GetObject(&p)

	case "DeleteObjectInput":
		p := input.(s3.DeleteObjectInput)
		output, err = service.DeleteObject(&p)
	}

	if err != nil {
		log.Println("S3Request Err:", err)
		return nil, err
	}
	return output, nil
}
