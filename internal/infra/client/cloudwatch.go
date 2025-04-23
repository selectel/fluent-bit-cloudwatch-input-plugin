package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/model"
)

type Cloudwatch struct {
	client *cloudwatchlogs.Client
}

func NewCloudwatchClient(ctx context.Context, region, endpoint string) (*Cloudwatch, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithBaseEndpoint(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	client := &Cloudwatch{
		client: cloudwatchlogs.NewFromConfig(cfg),
	}

	return client, nil
}

func (cw *Cloudwatch) GetLogEvents(
	ctx context.Context,
	logGroupName, logStreamName, nextToken string,
) ([]model.Event, string, error) {
	params := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(true),
	}

	if nextToken != "" {
		params.NextToken = aws.String(nextToken)
	}

	resp, err := cw.client.GetLogEvents(ctx, params)
	if err != nil {
		return nil, "", fmt.Errorf("failed to complete request: %w", err)
	}

	events := make([]model.Event, 0, len(resp.Events))

	for _, event := range resp.Events {
		events = append(events, model.Event{
			IngestionTime: *event.IngestionTime,
			Timestamp:     *event.Timestamp,
			Message:       *event.Message,
		})
	}

	return events, *resp.NextForwardToken, nil
}
