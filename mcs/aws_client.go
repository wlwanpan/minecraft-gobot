package mcs

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wlwanpan/minecraft-gobot/config"
)

var (
	ErrInstanceIdNotFound = errors.New("instance id not found")

	ErrNoRunningInstances = errors.New("no running instances")
)

type EC2StatusResp struct {
	StateCode instanceStateCode
	State     string
}

type S3StoreFileResp struct {
	Name  string
	Size  int64
	S3URL string
}

func NewSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(config.Cfg.Aws.Region),
	})
}

type AWSClient struct {
	session *session.Session
	svc     *ec2.EC2
	storage *s3.S3
}

func New() (*AWSClient, error) {
	s, err := NewSession()
	if err != nil {
		return nil, err
	}

	return &AWSClient{
		session: s,
		svc:     ec2.New(s),
		storage: s3.New(s),
	}, nil
}

func (c *AWSClient) StoreFile(zipPath string, filename string) (*S3StoreFileResp, error) {
	file, err := os.Open(zipPath)
	if err != nil {
		return nil, err
	}
	stat, _ := file.Stat()
	size := stat.Size()
	fileBuffer := make([]byte, size)
	file.Read(fileBuffer)

	log.Printf("aws s3: uploading file='%s'", stat.Name())

	s3BucketName := config.Cfg.Aws.S3BucketName
	s3BucketRegion := config.Cfg.Aws.Region
	output, err := c.storage.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3BucketName),
		Key:                  aws.String(filename),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("s3 PutObjectOutput: %s", output.String())

	return &S3StoreFileResp{
		Name:  filename,
		Size:  size,
		S3URL: fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3BucketRegion, s3BucketName, filename),
	}, nil
}
