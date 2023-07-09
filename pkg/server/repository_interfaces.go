package server

import "context"

type coverageValueGetter interface {
	GetValue1000(ctx context.Context, key string) (int, error)
}

type coverageValueStorer interface {
	StoreValue1000(ctx context.Context, key string, value1000 int) error
}

type CoverageRepository interface {
	coverageValueGetter
	coverageValueStorer
}
