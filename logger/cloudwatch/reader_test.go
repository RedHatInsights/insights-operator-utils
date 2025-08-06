package cloudwatch

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReader(t *testing.T) {
	c := new(mockClient)
	r := &Reader{
		group:  aws.String("group"),
		stream: aws.String("1234"),
		client: c,
	}

	c.On("GetLogEvents", mock.Anything, &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String("group"),
		StartFromHead: aws.Bool(true),
		LogStreamName: aws.String("1234"),
	}).Once().Return(&cloudwatchlogs.GetLogEventsOutput{
		Events: []types.OutputLogEvent{
			{Message: aws.String("Hello"), Timestamp: aws.Int64(1000)},
		},
	}, nil)

	err := r.read()
	assert.NoError(t, err)

	b := make([]byte, 1000)
	n, err := r.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	c.AssertExpectations(t)
}

func TestReader_Buffering(t *testing.T) {
	c := new(mockClient)
	r := &Reader{
		group:  aws.String("group"),
		stream: aws.String("1234"),
		client: c,
	}

	c.On("GetLogEvents", mock.Anything, &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String("group"),
		StartFromHead: aws.Bool(true),
		LogStreamName: aws.String("1234"),
	}).Once().Return(&cloudwatchlogs.GetLogEventsOutput{
		Events: []types.OutputLogEvent{
			{Message: aws.String("Hello"), Timestamp: aws.Int64(1000)},
		},
	}, nil)

	err := r.read()
	assert.NoError(t, err)

	b := make([]byte, 3)
	n, err := r.Read(b) //Hel
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = r.Read(b) //lo
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	c.AssertExpectations(t)
}

func TestReader_EndOfFile(t *testing.T) {
	c := new(mockClient)
	r := &Reader{
		group:  aws.String("group"),
		stream: aws.String("1234"),
		client: c,
	}

	c.On("GetLogEvents", mock.Anything, &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String("group"),
		StartFromHead: aws.Bool(true),
		LogStreamName: aws.String("1234"),
	}).Once().Return(&cloudwatchlogs.GetLogEventsOutput{
		Events: []types.OutputLogEvent{
			{Message: aws.String("Hello"), Timestamp: aws.Int64(1000)},
		},
		NextForwardToken: aws.String("next"),
	}, nil)

	c.On("GetLogEvents", mock.Anything, &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String("group"),
		LogStreamName: aws.String("1234"),
		StartFromHead: aws.Bool(true),
		NextToken:     aws.String("next"),
	}).Once().Return(&cloudwatchlogs.GetLogEventsOutput{
		Events: []types.OutputLogEvent{},
	}, nil)

	err := r.read()
	assert.NoError(t, err)

	b := make([]byte, 1000)
	n, err := r.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	err = r.read()
	assert.NoError(t, err)

	n, err = r.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	c.AssertExpectations(t)
}

func TestReader_Err(t *testing.T) {
	c := new(mockClient)
	r := &Reader{
		group:  aws.String("group"),
		stream: aws.String("1234"),
		client: c,
	}

	c.On("GetLogEvents", mock.Anything, mock.Anything).Return(&cloudwatchlogs.GetLogEventsOutput{}, errors.New("error"))

	err := r.read()
	assert.Error(t, err)

	// Set the error manually since we're not using the goroutine
	r.err = err

	b := make([]byte, 1000)
	_, err = r.Read(b)
	assert.Error(t, err)

	c.AssertExpectations(t)
}
