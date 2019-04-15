package main

import (
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

//S3Storage : implement DataStorage Interface using AWS S3 Service
type S3Storage struct {
	S3Service *s3.S3
	Bucket    string
	Region    string
}

// NewS3Storage : Initialize S3 service and set default values
func NewS3Storage(bucket, region string) (*S3Storage, error) {
	s3Storage := S3Storage{Bucket: bucket, Region: region}
	service, err := NewS3Service("", "default", "ap-northeast-1")
	if err != nil || nil == service {
		log.Println("S3Storage Create s3 service Fail")
		return nil, err
	}
	s3Storage.S3Service = service
	return &s3Storage, nil
}

//GetData : Implement DataStorage Interface method of GetData
func (s *S3Storage) GetData(name string) ([]byte, error) {
	param := s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(name),
	}
	resp, err := S3Request(s.S3Service, param)
	if err != nil {
		log.Println("GetData S3Request error:", err)
		return nil, err
	}
	object := resp.(*s3.GetObjectOutput)
	//For Debug: log.Println("GetData resp", object)
	fileBytes, err := ioutil.ReadAll(object.Body)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
