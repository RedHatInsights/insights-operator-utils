package cloudwatch

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// Throttling and limits from http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/cloudwatch_limits.html
const (
	// The maximum rate of a GetLogEvents request is 10 requests per second per AWS account.
	readThrottle = time.Second / 10

	// The maximum rate of a PutLogEvents request is 5 requests per second per log stream.
	writeThrottle = time.Second / 5

	// maximum message size is 1048576, but we have some metadata
	maximumBatchSize = 1048576 / 2
)

// now is a function that returns the current time.Time. It's a variable so that
// it can be stubbed out in unit tests.
var now = time.Now

// client duck types the aws sdk client for testing.
type client interface {
	PutLogEvents(context.Context, *cloudwatchlogs.PutLogEventsInput, ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogEventsOutput, error)
	CreateLogStream(context.Context, *cloudwatchlogs.CreateLogStreamInput, ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogStreamOutput, error)
	GetLogEvents(context.Context, *cloudwatchlogs.GetLogEventsInput, ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error)
	DescribeLogStreams(context.Context, *cloudwatchlogs.DescribeLogStreamsInput, ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
}

// Group wraps a log stream group and provides factory methods for creating
// readers and writers for streams.
type Group struct {
	group  string
	client *cloudwatchlogs.Client
}

// NewGroup returns a new Group instance.
func NewGroup(group string, client *cloudwatchlogs.Client) *Group {
	return &Group{
		group:  group,
		client: client,
	}
}

// Existing uses an existing log group created previously
func (g *Group) existing(stream string) (io.Writer, error) {
	result, err := g.client.DescribeLogStreams(context.TODO(), &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(g.group),
		LogStreamNamePrefix: aws.String(stream),
		OrderBy:             types.OrderByLogStreamName,
		Descending:          aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	if len(result.LogStreams) == 0 {
		return nil, errors.New("Log stream not found " + stream)
	}

	// since values are sorted the stream with exact match will be first
	logStream := result.LogStreams[0]

	return NewWriterWithToken(g.group, stream, logStream.UploadSequenceToken, g.client), nil
}

// Create creates a log stream in the group and returns an io.Writer for it.
func (g *Group) Create(stream string) (io.Writer, error) {
	if _, err := g.client.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(g.group),
		LogStreamName: aws.String(stream),
	}); err != nil {
		var resourceAlreadyExists *types.ResourceAlreadyExistsException
		if errors.As(err, &resourceAlreadyExists) {
			return g.existing(stream)
		}
		return nil, err
	}

	return NewWriter(g.group, stream, g.client), nil
}

// Open returns an io.Reader to read from the log stream.
func (g *Group) Open(stream string) (io.Reader, error) {
	return NewReader(g.group, stream, g.client), nil
}
