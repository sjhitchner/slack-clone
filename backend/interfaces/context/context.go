package context

import (
	"context"

	"github.com/sjhitchner/slack-clone/backend/domain"
)

type ContextKey int32

const (
	InteractorKey ContextKey = iota
	AggregatorKey
	UserIdKey
	IsAuthorizedKey
)

func SetAggregator(ctx context.Context, agg domain.Aggregator) context.Context {
	return context.WithValue(ctx, AggregatorKey, agg)
}

func Aggregator(ctx context.Context) domain.Aggregator {
	return ctx.Value(AggregatorKey).(domain.Aggregator)
}

func SetInteractor(ctx context.Context, inter domain.Interactor) context.Context {
	return context.WithValue(ctx, InteractorKey, inter)
}

func Interactor(ctx context.Context) domain.Interactor {
	return ctx.Value(InteractorKey).(domain.Interactor)
}

func SetCurrentUserId(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, UserIdKey, userId)
}

func CurrentUserId(ctx context.Context) int64 {
	return ctx.Value(UserIdKey).(int64)
}
