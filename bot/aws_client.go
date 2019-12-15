package bot

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	AWS_REGION string = "ca-central-1"

	EC2_INSTANCE_ID string = "i-0dd138d76b910e575"

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

type awsEC2StatusResp struct {
	StateCode instanceStateCode
	State     string
}

type awsClient struct {
	session *session.Session
	svc     *ec2.EC2
}

func NewAwsClient() (*awsClient, error) {
	s, err := session.NewSession(&aws.Config{Region: &AWS_REGION})
	if err != nil {
		return nil, err
	}

	return &awsClient{
		session: s,
		svc:     ec2.New(s),
	}, nil
}

func (c *awsClient) StartInstance() error {
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

func (c *awsClient) StopInstance() error {
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

func (c *awsClient) InstanceStatus() (*awsEC2StatusResp, error) {
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

	return &awsEC2StatusResp{
		StateCode: instanceStateCode(*status.InstanceState.Code),
		State:     *status.InstanceStatus.Status,
	}, nil
}
