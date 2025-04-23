package adapter

import "context"

type Storage interface {
	GetNextToken(ctx context.Context, region, logGroupName, logStreamName string) (string, error)
	SetNextToken(ctx context.Context, region, logGroupName, logStreamName, nextToken string) error
}
