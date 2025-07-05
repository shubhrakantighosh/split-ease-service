package service

import "context"

type Interface interface {
	Get(ctx context.Context)
}
