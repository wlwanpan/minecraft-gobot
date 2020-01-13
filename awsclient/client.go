package awsclient

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
)

var (
	AWS_REGION string = "ca-central-1" // TODO: Move to config file

	AWS_S3_BUCKET string = "elasticbeanstalk-ca-central-1-597927659010" // TODO: Move to config file

	EC2_INSTANCE_ID string = "i-0dd138d76b910e575" // TODO: Move to config file

	ErrInstanceIdNotFound = errors.New("instance id not found")

	ErrNoRunningInstances = errors.New("no running instances")
)

const (
	INSTANCE_PENDING_CODE       instanceStateCode = 0
	INSTANCE_RUNNING_CODE       instanceStateCode = 16
	INSTANCE_SHUTTING_DOWN_CODE instanceStateCode = 32
	INSTANCE_TERMINATED_CODE    instanceStateCode = 48
	INSTANCE_STOPPING_CODE      instanceStateCode = 64
	INSTANCE_STOPPED_CODE       instanceStateCode = 80
)

type instanceStateCode int

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
	return session.NewSession(&aws.Config{Region: &AWS_REGION})
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

	output, err := c.storage.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AWS_S3_BUCKET),
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
		S3URL: fmt.Sprintf("https://%s.s3.ca-central-1.amazonaws.com/%s", AWS_S3_BUCKET, filename),
	}, nil
}

func (c *AWSClient) StartInstance() error {
	config := &ec2.StartInstancesInput{
		InstanceIds: []*string{&EC2_INSTANCE_ID},
	}
	out, err := c.svc.StartInstances(config)
	if err != nil {
		return err
	}

	log.Printf("Starting instance: %s", out.GoString())
	return nil
}

func (c *AWSClient) StopInstance() error {
	config := &ec2.StopInstancesInput{
		InstanceIds: []*string{&EC2_INSTANCE_ID},
	}
	out, err := c.svc.StopInstances(config)
	if err != nil {
		return err
	}

	log.Printf("Stopping instance: %s", out.GoString())
	return nil
}

func (c *AWSClient) InstanceStatus() (*EC2StatusResp, error) {
	config := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&EC2_INSTANCE_ID},
	}
	out, err := c.svc.DescribeInstanceStatus(config)
	if err != nil {
		return nil, err
	}

	log.Printf("Raw response: %s", out.GoString())

	if len(out.InstanceStatuses) == 0 {
		return nil, ErrNoRunningInstances
	}

	status := out.InstanceStatuses[0]

	return &EC2StatusResp{
		StateCode: instanceStateCode(*status.InstanceState.Code),
		State:     *status.InstanceStatus.Status,
	}, nil
}
